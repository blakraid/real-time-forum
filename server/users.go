package server

import (
	 "database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"rtf/database"
)

type user struct{
	Username string
	Date sql.NullString
}
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err != nil {
		//http.Error(w, "You are not logged in", http.StatusBadRequest)
		response := map[string]string{"message": "You are not logged in"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var allUsers []user
	currentUsername := r.URL.Query().Get("username")

	// Updated SQL: sort by last message timestamp, fallback to username
	rows, err := database.Sql.Query(`
	SELECT
		u.username,
		m.timestamp
	FROM users u
	LEFT JOIN (
		SELECT 
			CASE 
				WHEN sender = ? THEN receiver
				ELSE sender
			END AS other_user,
			MAX(timestamp) AS timestamp
		FROM messages
		WHERE sender = ? OR receiver = ?
		GROUP BY other_user
	) m ON u.username = m.other_user
	WHERE u.username != ?;
	`,currentUsername,currentUsername,currentUsername,currentUsername,)

	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		fmt.Println("Query error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data user
		if err := rows.Scan(&data.Username ,&data.Date); err != nil {
			http.Error(w, "Database scan error", http.StatusInternalServerError)
			fmt.Println("Scan error:", err)
			return
		}
		if data.Username == currentUsername {
			continue 
		}
		allUsers = append(allUsers, data)
	}

	// Optionally log or debug
	//fmt.Printf("User %s sees: %v\n", currentUsername, allUsers)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allUsers); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		fmt.Println("Encoding error:", err)
	}

}
