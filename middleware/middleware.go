package middleware

import (
	"net/http"
	"fmt"
	"encoding/json"
	
	"example.com/server/types"
	"example.com/server/users"
	"example.com/server/database"
)

type AuthHandler struct {
	db *database.Database
	logger types.Logger
}

// Simple logger implementation
type ConsoleLogger struct{}

func NewAuthHandler(db *database.Database, logger types.Logger) *AuthHandler {
	return &AuthHandler {
		db: db,
		logger: logger,
	}
}

func (l *ConsoleLogger) Info(message string) {
	fmt.Println("INFO:", message)
}

func (l *ConsoleLogger) Error(message string) {
	fmt.Println("ERROR:", message)
}

func (h *AuthHandler) Signin(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Attempting to sign up a new user")

	var creds users.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		h.logger.Error("Error parsing request body: " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the parsed data
	if creds.Username == "" || creds.Password == "" {
		h.logger.Error("Validation failed: username or password is empty")
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	_, err = h.db.Query("INSERT INTO users (username, password) VALUES (?, ?)", creds.Username, creds.Password)
	if err != nil {
		h.logger.Error("Failed to insert new user")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	h.logger.Info("User signed up successfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Signin successful"))
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("attempting to sign in user")

	var creds users.Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil{
		h.logger.Error("Validation failed: username or password is empty")
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	h.logger.Info("User signed in successfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Signin successful"))
}