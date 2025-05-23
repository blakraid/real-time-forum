package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"rtf/database"
	"time"
)

func CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
	username, _, loggedIn, err := RequireLogin(w, r)
	if err != nil {
		fmt.Println("Error in the RequiredLogin !!! :", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if loggedIn {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"loggedIn": true,"username": "%s"}`, username)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"loggedIn": false}`)
	}
}

func RequireLogin(w http.ResponseWriter, r *http.Request) (string, string, bool, error) {
	cookie, _ := r.Cookie("session_token")
	if cookie == nil {
		return "", "guest", false, nil
	}

	var username, sessionToken string
	err := database.Sql.QueryRow("SELECT username, session_token FROM users WHERE session_token = ?", cookie.Value).Scan(&username, &sessionToken)
	if err == sql.ErrNoRows {
		cookies := r.Cookies()
		// Loop through the cookies and expire them
		for _, cookie := range cookies {
			http.SetCookie(w, &http.Cookie{
				Name:    cookie.Name,
				Value:   "",
				Expires: time.Now().Add(-1 * time.Hour),
				MaxAge:  -1,
			})
		}
		return "", "guest", false, err
	} else if err != nil {
		log.Printf("Database error: %v", err)
		return "", "guest", false, err
	}

	return username, sessionToken, true, nil
}
