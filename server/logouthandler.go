package server

import (
	"encoding/json"
	"net/http"
	"time"
)
var sessionStore = make(map[string]string)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		//http.Error(w, "You are not logged in", http.StatusBadRequest)
		response := map[string]string{"message": "You are not logged in"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Remove the session from the session store
	delete(sessionStore, cookie.Value)

	// Expire the cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "guest",
		Expires: time.Now().Add(-1 * time.Hour),
	})
}
