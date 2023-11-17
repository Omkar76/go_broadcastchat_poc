package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	_ "github.com/lib/pq"
	"nhooyr.io/websocket"
)

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type Manager struct {
	clients map[*Client]bool
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[*Client]bool),
	}
}

func (m *Manager) Add(c *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[c] = true
	go c.processReads()
	fmt.Println("Adding...Current client count ", len(m.clients))
}

type Client struct {
	Conn    *websocket.Conn
	manager *Manager
}

func (c *Client) processReads() {
	for {
		fmt.Println("Reading...")
		kind, data, err := c.Conn.Read(context.Background())
		fmt.Println(kind, data, err)

		if err != nil {
			log.Println(err)
			c.manager.Remove(c)
			break
		}
		c.manager.Broadcast(data)
		log.Println(kind, string(data))
	}
}

func (m *Manager) Remove(c *Client) {
	m.Lock()
	defer m.Unlock()
	delete(m.clients, c)

	c.Conn.Close(websocket.StatusNormalClosure, "Closing connection")
	fmt.Println("Removing...Current client count ", len(m.clients))
}

func (m *Manager) Broadcast(msg []byte) {
	fmt.Println("Broadcasting")
	m.Lock()
	defer m.Unlock()
	for c := range m.clients {
		fmt.Print("Sending to client")
		c.Conn.Write(context.Background(), websocket.MessageText, msg)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	var c *websocket.Conn
	var err error

	manager := NewManager()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err = websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
			CompressionMode:    websocket.CompressionDisabled,
		})

		fmt.Println("newClient connectesd")
		if err != nil {
			log.Println(err)
			return
		}
		client := &Client{
			Conn:    c,
			manager: manager,
		}
		manager.Add(client)
	})

	http.ListenAndServe(":8000", nil)
}
