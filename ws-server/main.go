package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	ALPHANUMERIC_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	DEFAULT_PAYLOAD_SIZE_KB = 1
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Untuk development saja, di production harus dibatasi
	},
}

func generateRandomPayload(sizeKB int) string {
	totalChars := sizeKB * 1024
	var builder strings.Builder
	builder.Grow(totalChars)

	for i := 0; i < totalChars; i++ {
		randomIndex := rand.Intn(len(ALPHANUMERIC_CHARS))
		builder.WriteByte(ALPHANUMERIC_CHARS[randomIndex])
	}

	return builder.String()
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Get payload size from query parameter
	sizeKB := DEFAULT_PAYLOAD_SIZE_KB
	if sizeParam := r.URL.Query().Get("size"); sizeParam != "" {
		if size, err := strconv.Atoi(sizeParam); err == nil && size > 0 {
			sizeKB = size
		}
	}

	// Generate payload
	payload := generateRandomPayload(sizeKB)
	timestamp := time.Now().Format(time.RFC3339Nano)

	// Create response message
	response := map[string]interface{}{
		"message":       "Hello from WebSocket backend!",
		"timestamp":     timestamp,
		"payload_size":  sizeKB,
		"payload":       payload,
	}

	// Send response
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Write error: %v", err)
		return
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	http.HandleFunc("/ws", wsHandler)
	log.Printf("WebSocket server running on :9002")
	log.Printf("Default payload size: %d KB", DEFAULT_PAYLOAD_SIZE_KB)
	log.Println("Use ?size=X query parameter to specify payload size in KB")
	
	log.Fatal(http.ListenAndServe(":9002", nil))
}