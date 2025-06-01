package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"rtf/database"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// User struct
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Message struct
type Message struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

var (
	clients   = make(map[string][]*websocket.Conn)
	clientsMu sync.Mutex
)

// Handle WebSocket connection
func HandleConnection(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Add user to online list
	clientsMu.Lock()
	clients[userID] = append(clients[userID], conn)
	clientsMu.Unlock()
	

	BroadcastUserList() // Notify all clients of new user

	// fmt.Println("User connected:", userID)

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			// fmt.Println("Error reading message:", err)
			break
		}

		if msg.Type == "message" {
			// Save message to DB
			_, err = database.Sql.Exec("INSERT INTO messages (sender, receiver, text) VALUES (?, ?, ?)", msg.Sender, msg.Receiver, msg.Text)
			if err != nil {
				fmt.Println("Error saving message:", err)
				continue
			}

			clientsMu.Lock()
			if receiverConns, ok := clients[msg.Receiver]; ok {
				for _, conn := range receiverConns {
					err := conn.WriteJSON(msg)
					if err != nil {
						fmt.Println("Error sending to receiver:", err)
					}
				}
			}
			clientsMu.Unlock()
			clientsMu.Lock()
			msg.Type = "msg"
			if my, ok := clients[userID]; ok {
				for _, conn := range my {
					err := conn.WriteJSON(msg)
					if err != nil {
						fmt.Println("Error sending to receiver:", err)
					}
				}
			}
			clientsMu.Unlock()
		} else if msg.Type == "typing" {
			// Prepare JSON response
			usersJSON1, _ := json.Marshal(map[string]interface{}{
				"type":  "typing",
				"text":  "typing",
				"users": userID,
			})
			for _, conn := range clients[msg.Receiver] {
				conn.WriteMessage(websocket.TextMessage, usersJSON1)
			}
		}
	}

	// Remove user on disconnect
	clientsMu.Lock()
	delete(clients, userID)
	clientsMu.Unlock()

	BroadcastUserList() // Notify all users of updated list
}

// Fetch last 10 messages
func GetMessages(w http.ResponseWriter, r *http.Request) {
	sender := r.URL.Query().Get("sender")
	receiver := r.URL.Query().Get("receiver")
	num := r.URL.Query().Get("num")

	rows, err := database.Sql.Query(`
	SELECT sender, receiver, text, timestamp FROM messages
	WHERE (sender=? AND receiver=?) OR (sender=? AND receiver=?)
	ORDER BY timestamp DESC LIMIT 10 OFFSET ?`, sender, receiver, receiver, sender, num)
	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}
	// fmt.Println(sender + "  " + receiver)
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		rows.Scan(&msg.Sender, &msg.Receiver, &msg.Text, &msg.Timestamp)
		messages = append(messages, msg)
	}

	json.NewEncoder(w).Encode(messages)
}

func BroadcastUserList() {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	// Get list of online usernames
	var userList []string
	for username := range clients {
		userList = append(userList, username)
	}

	// Prepare JSON response
	usersJSON, _ := json.Marshal(map[string]interface{}{
		"type":  "userList",
		"users": userList,
	})

	// Send to all connected clients
	for _, conn := range clients {
		for _, conn1 := range conn {
			conn1.WriteMessage(websocket.TextMessage, usersJSON)
		}
	}
}


// func BroadcastMySelf(msg Message){
// 	usersJSON, _ := json.Marshal(map[string]interface{}{
// 		"type":  "myWebSocket",
// 		"msg": msg,
// 	})
// 		for _, conn1 := range clients[msg.Sender] {
// 			conn1.WriteMessage(websocket.TextMessage, usersJSON)
// 		}
// }