package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rtf/packages"
	"time"
	"rtf/database"
	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	email := tools.EscapeString(r.FormValue("email"))
	password := tools.EscapeString(r.FormValue("password"))
	

	const maxEmail = 100
	const maxPassword = 100

	if len(email) > maxEmail {
		response := map[string]string{"message": fmt.Sprintf("Email cannot be longer than %d characters.", maxEmail)}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(password) > maxPassword {
		response := map[string]string{"message": fmt.Sprintf("Password cannot be longer than %d characters.", maxPassword)}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var storedPassword, sessionToken, username string
	err := database.Sql.QueryRow("SELECT password, session_token, username FROM users WHERE email = ? OR username = ?", email, email).Scan(&storedPassword, &sessionToken, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]string{"message": "Invalid info"}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
		} else {
			log.Printf("Database error: %v", err)
			response := map[string]string{"message": "Internal server error"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		response := map[string]string{"message": "Invalid info"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate a new session token
	newSessionToken, _ := uuid.NewV4()
	sessionToken = newSessionToken.String()

	// Update the session token in the database
	_, err = database.Sql.Exec("UPDATE users SET session_token = ? WHERE email = ? or username = ?", sessionToken, email,email)
	if err != nil {
		log.Printf("Error updating session token: %v", err)
		response := map[string]string{"error": "Internal server error"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Set a cookie with the session token
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(1 * time.Hour),
	})

	response := map[string]string{"message": "Login successful!", "username": username}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
