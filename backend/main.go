package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     int = 4
	MaxCost     int = 31
	DefaultCost int = 10
)

var validate *validator.Validate

// Using pointer and "omitempty" in json tag for optinal tag
// Because if the value is missing in the unmarshal step
// the value will record as empty
// by using * and "omitempty" there no field in the json object
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int       `bun:",pk,autoincrement"   json:"id"`
	Name          string    `bun:",unique,notnull"     json:"name"       validate:"required,alphanumunicode"`
	Email         string    `bun:",unique,notnull"     json:"email"      validate:"required,email"`
	Password      string    `bun:",notnull"            json:"password"   validate:"required,alphanumunicode"`
	LastLogin     time.Time `bun:""                    json:"last_login"                                     vaidate:"reuired,datetime"`
}

type LoginForm struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,alphanumunicode"`
}

type RegisterForm struct {
	Name     string `json:"name"     validate:"required,alphanumunicode"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,alphanumunicode"`
}

type Claims struct {
	Email              string `json:"email" validate:"required,email"`
	jwt.StandardClaims `                    validate:"required"`
}

// Get local map secrets
func getEnvMap() (localEnv map[string]string, err error) {
	localEnv, err = godotenv.Read()
	return
}

// Connect to Postgres database
func connectDB() (db *bun.DB) {
	localEnv, err := getEnvMap()
	if err != nil {
		log.Panic("Failed to load .env file")
	}

	dsn := localEnv["DB_POSTGRES"]
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(pgdb, pgdialect.New())
	return
}

// Encrypt user password
func hashPassword(password string) (string, error) {
	// Hashing password byte with default cost
	hashedPassByte, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	return string(hashedPassByte), err
}

// Match user password
func matchPassword(inputPass string, dbPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(inputPass), []byte(inputPass))
	return err == nil
}

func generateToken(email string) (tokenString string, err error) {
	localEnv, err := getEnvMap()
	if err != nil {
		log.Panic("Failed to load .env file")
	}

	secretString := localEnv["JWT_SECRET"]
	secretByte := []byte(secretString)

	expirationTimes := time.Now().AddDate(0, 0, 7)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimes.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(secretByte)
	return
}

// Create reset users table
func createUsersTable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	db := connectDB()

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	if err := db.ResetModel(ctx, (*User)(nil)); err != nil {
		log.Panic(err)
	}

	w.WriteHeader(http.StatusCreated)
}

// Generate user session
func generateSession(w http.ResponseWriter, r *http.Request, userEmail string) {
	tokenString, err := generateToken(userEmail)
	var headerToken string = "Bearer " + tokenString

	if err != nil {
		log.Panic("Failed to get generated token string")
	}

	// Write token to return header
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Authorization", headerToken)
	http.Redirect(w, r, "localhost:3000", http.StatusOK)
}

// Add new register user to database
// Get user information from register form
// Then unmarshal and create new user object
// Then add object to database
func addNew(w http.ResponseWriter, r *http.Request) {
	var newRegister RegisterForm
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request mising values"))
	}

	json.Unmarshal(reqBody, &newRegister)

	// Register form validation
	validate = validator.New()
	if err := validate.Struct(newRegister); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Panic(err)
		}

		var errListMes string

		for _, err := range err.(validator.ValidationErrors) {
			errListMes += ", " + err.StructField()
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing values:" + errListMes[2:]))
	}

	// Prepare database
	ctx := context.Background()
	db := connectDB()

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	// Check duplidate email in database
	exists, err := db.NewSelect().
		Model((*User)(nil)).
		Where("email = ?", newRegister.Email).
		WhereOr("name = ?", newRegister.Name).
		Exists(ctx)
	if err != nil {
		log.Panic("Failed to check user exists")
	}

	if exists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User exists"))
	}

	// Encrypt user password
	hashedPass, err := hashPassword(newRegister.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to hashing user password"))
	}

	newUser := &User{
		Name:      newRegister.Name,
		Email:     newRegister.Email,
		Password:  hashedPass,
		LastLogin: time.Now(),
	}

	// Save new user in database
	if _, err := db.NewInsert().Model(newUser).Exec(ctx); err != nil {
		log.Panic("Failed to save new user")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Create login token
	generateSession(w, r, newRegister.Email)
}

// Validate logine password and the hashed password
// based on the user email
func validateLogin(w http.ResponseWriter, r *http.Request) {
	var newLogin LoginForm
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing values from the request body"))
	}

	json.Unmarshal(reqBody, &newLogin)

	// New login validation
	validate = validator.New()
	if err := validate.Struct(newLogin); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Panic(err)
		}

		var errListMes string

		for _, err := range err.(validator.ValidationErrors) {
			errListMes += ", " + err.StructField()
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing values: " + errListMes))
	}

	// Loooking for user email in database
	ctx := context.Background()
	db := connectDB()

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	foundUser := new(User)

	err = db.NewSelect().Model(foundUser).Where("email = ?", newLogin.Email).Scan(ctx)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User not found on database"))
	}

	// Matching password
	if isMatch := matchPassword(newLogin.Password, foundUser.Password); !isMatch {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Password is not correct"))
	}

	// Update user last login session
	foundUser.LastLogin = time.Now()
	if _, err = db.NewUpdate().Model(foundUser).WherePK().Exec(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed udpate user last login session"))
	}

	// Generate new session
	generateSession(w, r, newLogin.Email)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome homepage"))
}

func startRouter() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/users/table", createUsersTable).Methods("POST")
	r.HandleFunc("/auth/register", addNew).Methods("POST")
	r.HandleFunc("/auth/login", validateLogin).Methods("POST")

	// Start server
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:3001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func main() {
	log.Println("Starting server...")
	startRouter()
}
