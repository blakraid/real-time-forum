package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"rtf/database"
	tools "rtf/packages"
	"strconv"
	"strings"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegesterHandler(w http.ResponseWriter, r *http.Request) {
	
	var response = make(map[string]string)

	if r.Method != http.MethodPost {
		response = map[string]string{"message": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	username := tools.EscapeString(r.FormValue("username"))
	FirtsName := strings.ToLower(tools.EscapeString(r.FormValue("FirtsName")))
	LastName := strings.ToLower(tools.EscapeString(r.FormValue("LastName")))
	age , err := strconv.Atoi(tools.EscapeString(r.FormValue("age")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"message": "invalid age"}
		json.NewEncoder(w).Encode(response)
		return
	}
	gender := strings.ToLower(tools.EscapeString(r.FormValue("gender")))

	err1 := validotherdata(FirtsName,LastName,gender)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"message": err1.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}


	email := tools.EscapeString(r.FormValue("email"))
	password := tools.EscapeString(r.FormValue("password"))

	_, valid := ValidateInput(username, email, password)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{"message": "Password or mail is wrong"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var existingUsername, existingEmail string
	err = database.Sql.QueryRow("SELECT username, email FROM users WHERE email = ? OR username = ?", email, username).Scan(&existingUsername, &existingEmail)
	if err == nil {
		var conflictField, conflictMessage string
		if existingUsername == username {
			conflictField = "username"
			conflictMessage = "Username already exists"
		} else if existingEmail == email {
			conflictField = "email"
			conflictMessage = "Email already exists"
		}

		response = map[string]string{"message": conflictMessage,"message1":conflictField}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	} else if err != sql.ErrNoRows {
		log.Printf("Database error: %v", err)
		response = map[string]string{"message": "Database error"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response = map[string]string{"message": "Error hashing password"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	sessionToken, _ := uuid.NewV4()

	_, err = database.Sql.Exec("INSERT INTO users (username, email, password,firtsName, lastName, age,gender, session_token) VALUES (?, ?, ?, ?,?, ?, ?, ?)", username, email, hashedPassword,FirtsName, LastName, age,gender, sessionToken)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		response = map[string]string{"message": "Registration failed"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response = map[string]string{"message": "Registration successful!"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ValidateInput(username, email, password string) (map[string]string, bool) {
	errors := make(map[string]string)

	const maxUsername = 50
	const maxEmail = 100
	const maxPassword = 100

	if len(username) == 0 {
		errors["username"] = "Username cannot be empty"
		return errors, false
	} else if len(username) > maxUsername {
		errors["username"] = fmt.Sprintf("Username cannot be longer than %d characters.", maxUsername)
		return errors, false
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if len(email) == 0 {
		errors["email"] = "Email cannot be empty"
		return errors, false
	} else if len(email) > maxEmail {
		errors["email"] = fmt.Sprintf("Email cannot be longer than %d characters.", maxEmail)
		return errors, false
	} else if !emailRegex.MatchString(email) {
		errors["email"] = "Invalid email format"
		return errors, false
	}

	if len(password) < 8 {
		errors["password"] = "Password must be at least 8 characters long"
		return errors, false
	} else if len(password) > maxPassword {
		errors["password"] = fmt.Sprintf("Password cannot be longer than %d characters.", maxPassword)
		return errors, false
	} else {
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)

		if !hasUpper {
			errors["password"] = "Password must include at least one uppercase letter"
			return errors, false
		} else if !hasLower {
			errors["password"] = "Password must include at least one lowercase letter"
			return errors, false
		} else if !hasDigit {
			errors["password"] = "Password must include at least one digit"
			return errors, false
		} else if !hasSpecial {
			errors["password"] = "Password must include at least one special character"
			return errors, false
		}
	}

	if len(errors) > 0 {
		log.Println(errors)
		return errors, false
	}
	return nil, true
}


func validotherdata(firstname , lastname,gander string) error{
	firstregex := regexp.MustCompile(`^[a-z]+$`).MatchString(firstname)
	lastregex := regexp.MustCompile(`^[a-z]+$`).MatchString(lastname)
	if !firstregex || !lastregex {
		return errors.New("invalid first or last name")
	}
	if len(firstname) > 50 || len(lastname) > 50{
		return errors.New("first or last name is too long")
	}
	if gander != "man" && gander != "women" {
		return errors.New("invalid gander")
	}
	return nil
}