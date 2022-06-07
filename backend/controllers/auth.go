package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"rga/backend/config"
	"rga/backend/models"
	"rga/backend/utils"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Add new register user to database
// Get user information from register form
// Then unmarshal and create new user object
// Then add object to database
func Register(w http.ResponseWriter, r *http.Request) {
	var newRegister models.RegisterForm
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
	db := config.InitializeDB()

	// Check duplidate email in database
	exists, err := db.NewSelect().
		Model((*models.User)(nil)).
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
	hashedPass, err := utils.HashingPassword(newRegister.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to hashing user password"))
	}

	newUser := &models.User{
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
	utils.GenerateSession(w, r, newRegister.Email)
}

// Validate logine password and the hashed password
// based on the user email
func Login(w http.ResponseWriter, r *http.Request) {
	var newLogin models.LoginForm
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
	db := config.InitializeDB()

	foundUser := new(models.User)

	err = db.NewSelect().Model(foundUser).Where("email = ?", newLogin.Email).Scan(ctx)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User not found on database"))
	}

	// Matching password
	if isMatch := utils.MatchingPassword(newLogin.Password, foundUser.Password); !isMatch {
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
	utils.GenerateSession(w, r, newLogin.Email)
}
