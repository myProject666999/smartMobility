package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"smartMobility/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
}

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
	UserID  uint        `json:"user_id,omitempty"`
}

type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Mutex      sync.RWMutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var WSHub = &Hub{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[string]*Client),
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client.ID] = client
			h.Mutex.Unlock()
			log.Printf("WebSocket客户端连接: %s", client.ID)
		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
				log.Printf("WebSocket客户端断开: %s", client.ID)
			}
			h.Mutex.Unlock()
		case message := <-h.Broadcast:
			h.Mutex.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.ID)
				}
			}
			h.Mutex.RUnlock()
		}
	}
}

func (h *Hub) SendToUser(userID uint, message *Message) {
	h.Mutex.RLock()
	defer h.Mutex.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return
	}

	for _, client := range h.Clients {
		if client.UserID == userID {
			select {
			case client.Send <- messageBytes:
			default:
				close(client.Send)
				delete(h.Clients, client.ID)
			}
		}
	}
}

func (h *Hub) BroadcastToAll(message *Message) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return
	}
	h.Broadcast <- messageBytes
}

func (c *Client) Read() {
	defer func() {
		WSHub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(int64(config.AppConfig.Websocket.ReadBufferSize))

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket读取错误: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "ping":
			pongMsg := &Message{
				Type:    "pong",
				Content: "pong",
			}
			pongBytes, _ := json.Marshal(pongMsg)
			c.Send <- pongBytes
		case "message":
			msg.UserID = c.UserID
			WSHub.BroadcastToAll(&msg)
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	userID, _ := c.Get("user_id")
	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = conn.RemoteAddr().String()
	}

	client := &Client{
		ID:     clientID,
		UserID: userID.(uint),
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	WSHub.Register <- client

	go client.Write()
	go client.Read()
}

func SendNotification(userID uint, title, content string, notifType int) {
	message := &Message{
		Type: "notification",
		Content: map[string]interface{}{
			"title":   title,
			"content": content,
			"type":    notifType,
		},
	}
	WSHub.SendToUser(userID, message)
}
