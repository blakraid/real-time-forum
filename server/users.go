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
	Content sql.NullString
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
	m.timestamp AS last_msg_time,
	m.text AS last_msg_content
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
	) lm ON u.username = lm.other_user
	LEFT JOIN messages m ON 
	((m.sender = u.username AND m.receiver = ?) OR
	(m.receiver = u.username AND m.sender = ?))
	AND m.timestamp = lm.timestamp
	WHERE u.username != ?
	ORDER BY
	CASE WHEN lm.timestamp IS NULL THEN 1 ELSE 0 END,
	lm.timestamp DESC,
	u.username ASC;


	`,currentUsername,currentUsername,currentUsername,currentUsername,currentUsername,currentUsername)

	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		fmt.Println("Query error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data user
		if err := rows.Scan(&data.Username ,&data.Date,&data.Content); err != nil {
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
