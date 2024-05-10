package users

import (
	"database/sql"
	"log"
	"net/http"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"

	"example.com/server/database"
	"example.com/server/types"
)

type User struct {
	Username string
	Password string
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Response struct {
	
}

func RegisterUser(db *database.Database, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (username, password_hash, email) VALUES (?, ?, ?)", username, hashedPassword)
	return err
}

func AuthenticateUser(db *sql.DB, username, password string) bool {
	var hashedPassword string

	err := db.QueryRow("SELECT password_hash FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		log.Println("Authentication failed:", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.GetDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	rows, err := db.Query("SELECT id, username FROM users")
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Username); err != nil {
			http.Error(w, "Row scan error", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	response := types.Response {
		Status: "success",
		Message: "Users fetched successfully",
		Data: users,
	}

	w.Header().Set("Content-Type", "applcation/json")
	json.NewEncoder(w).Encode(response)
}
