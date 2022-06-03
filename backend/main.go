package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

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

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int       `bun:",pk,autoincrement"   json:"id"`
	Name          string    `bun:",notnull"            json:"name"`
	Email         string    `bun:",notnull"            json:"email"`
	Password      string    `bun:",notnull"            json:"password"`
	LastLogin     time.Time `bun:""                    json:"last_login"`
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	ConfirmPass string `json:"confirm_pass"`
}

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
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

func generateToken(email string, id string) (tokenString string, err error) {
	localEnv, err := getEnvMap()
	if err != nil {
		log.Panic("Failed to load .env file")
	}

	secretString := localEnv["JWT_SECRET"]
	secretByte := []byte(secretString)

	expirationTimes := time.Now().AddDate(0, 0, 7)
	claims := &Claims{
		ID:    id,
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

	// Save new user
	ctx := context.Background()
	db := connectDB()

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	res, err := db.NewInsert().Model(newUser).Returning("id").Exec(ctx)
	if err != nil {
		log.Panic("Failed to save new user")
	}

	// Create login token
	tokenString, err := generateToken(newRegister.Email, res)
	var headerToken string = "Bearer " + tokenString

	if err != nil {
		log.Panic("Failed to get generated token string")
	}

	// Write token to return header
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Authorization", headerToken)
	http.Redirect(w, r, "localhost:3000/home", http.StatusOK)
}

func validateLogin(w http.ResponseWriter, r *http.Request) {
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome homepage"))
}

func startRouter() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/users/table", createUsersTable).Methods("POST")

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
