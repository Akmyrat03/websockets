package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader to upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (not safe for production!)
	},
}

// WebSocket connections for the two users
var user1Conn *websocket.Conn
var user2Conn *websocket.Conn

// WebSocket handler
func chatHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}

	defer conn.Close()

	// Assign connection to either user1 or user2
	if user1Conn == nil {
		user1Conn = conn
		fmt.Println("User1 connected!")
	} else if user2Conn == nil {
		user2Conn = conn
		fmt.Println("User2 connected!")
	} else {
		conn.WriteMessage(websocket.TextMessage, []byte("Chat room is full!"))
		conn.Close()
		return
	}

	// Message listening loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Connection closed:", err)
			break
		}

		// Send the message to the other user
		if conn == user1Conn && user2Conn != nil {
			user2Conn.WriteMessage(websocket.TextMessage, message)
		} else if conn == user2Conn && user1Conn != nil {
			user1Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func main() {
	r := gin.Default()

	// Serve the WebSocket endpoint
	r.GET("/ws", chatHandler)

	// Serve the chat client HTML
	r.StaticFile("/", "./index.html")

	r.Run(":8080")
}
