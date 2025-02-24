package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	Conn   *websocket.Conn
	SendCh chan []byte
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			fmt.Println("Client registered:", client.ID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.SendCh)
				fmt.Println("Client unregistered:", client.ID)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.Lock()
			for _, client := range h.Clients {
				fmt.Println("Broadcasting message to client:", client.ID)
				client.SendCh <- message
			}
			h.mu.Unlock()
		}
	}
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	// for {
	// 	if false {
	// 		fmt.Print("")
	// 	}
	// }

	for message := range c.SendCh {
		fmt.Println("in write pump")

		var msg struct {
			SubmissionID string `json:"submission_id"`
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			panic(err)
		}

		if c.ID == msg.SubmissionID {
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				panic(err)
			}
		}

	}
}
