package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"rtf/database"
)

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err != nil {
		//http.Error(w, "You are not logged in", http.StatusBadRequest)
		response := map[string]string{"message": "You are not logged in"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var allUsers []string
	currentUsername := r.URL.Query().Get("username")

	// Updated SQL: sort by last message timestamp, fallback to username
	rows, err := database.Sql.Query(`
	SELECT u.username, MAX(m.timestamp) AS last_msg_time
	FROM users u
	LEFT JOIN (
		SELECT sender AS username, timestamp FROM messages
		UNION ALL
		SELECT receiver AS username, timestamp FROM messages
	) m ON u.username = m.username
	WHERE u.username != ?
	GROUP BY u.username
	ORDER BY 
		CASE WHEN last_msg_time IS NULL THEN 1 ELSE 0 END,
		last_msg_time DESC
`, currentUsername)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		fmt.Println("Query error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		var lastMsgTime sql.NullString
		if err := rows.Scan(&username, &lastMsgTime); err != nil {
			http.Error(w, "Database scan error", http.StatusInternalServerError)
			fmt.Println("Scan error:", err)
			return
		}
		if username == currentUsername {
			continue 
		}
		allUsers = append(allUsers, username)
	}

	// Optionally log or debug
	//fmt.Printf("User %s sees: %v\n", currentUsername, allUsers)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allUsers); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		fmt.Println("Encoding error:", err)
	}

}
