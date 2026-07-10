package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"turtle/core/lgr"
	"turtle/tools"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512 * 1024
)

// ─────────────────────────────────────────────
//  Myio hub
// ─────────────────────────────────────────────

type Myio struct {
	clients map[string]*Client
	mu      *tools.TimedMutex
}

func NewMyio() *Myio {
	tmp := &Myio{
		clients: make(map[string]*Client),
		mu:      tools.NewTimedMutex("MyioMutex", 5),
	}
	tmp.mu.Monitor(3 * time.Second)
	return tmp
}

func (self *Myio) addClient(c *Client) {
	self.mu.Lock("addClient")
	self.clients[c.uid] = c
	self.mu.Unlock()
}

func (self *Myio) removeClient(c *Client) {
	self.mu.Lock("removeClient")
	delete(self.clients, c.uid)
	self.mu.Unlock()
}

func (self *Myio) snapshot() map[string]*Client {
	self.mu.Lock("snapshot")
	defer self.mu.Unlock()
	out := make(map[string]*Client, len(self.clients))
	for k, v := range self.clients {
		out[k] = v
	}
	return out
}

// NotifyAll broadcasts data on the "notifications" channel to every client.
func (self *Myio) NotifyAll(data any) {
	self.Emit("notifications", data)
}

// Emit broadcasts to all connected clients.
func (self *Myio) Emit(channel string, data any) {
	defer tools.Recover("Failed to emit MYIO")
	for _, client := range self.snapshot() {
		client.emit(channel, data, nil)
	}
}

// EmitSync broadcasts synchronously to all connected clients.
func (self *Myio) EmitSync(channel string, data any) {
	defer tools.Recover("Failed to emit MYIO sync")
	for _, client := range self.snapshot() {
		client.emitSync(channel, data, nil)
	}
}

// EmitToRoom broadcasts to clients that have joined the given room.
func (self *Myio) EmitToRoom(room string, channel string, data any) {
	defer tools.Recover("Failed to emit MYIO room")
	for _, client := range self.snapshot() {
		client.mu.Lock("EmitToRoom-read")
		_, inRoom := client.rooms[room]
		client.mu.Unlock()
		if inRoom {
			client.emit(channel, data, nil)
		}
	}
}

// EmitToSession sends to all clients whose sessionId matches.
func (self *Myio) EmitToSession(sessionId string, channel string, data any) {
	self.EmitToRoom("session:"+sessionId, channel, data)
}

// EmitToUser sends to all clients whose userUid matches.
func (self *Myio) EmitToUser(userUid string, channel string, data any) {
	self.EmitToRoom("user:"+userUid, channel, data)
}

// ─────────────────────────────────────────────
//  Message protocol
// ─────────────────────────────────────────────

// Message is the wire format for all WebSocket frames.
type Message struct {
	Event   string      `json:"event"`
	Data    interface{} `json:"data"`
	AckId   *int        `json:"ackId,omitempty"`
	AckData interface{} `json:"ackData,omitempty"`
}

// ─────────────────────────────────────────────
//  Client
// ─────────────────────────────────────────────

type Client struct {
	uid   string
	conn  *websocket.Conn
	send  chan []byte
	mu    *tools.TimedMutex
	rooms map[string]struct{}
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		uid:   tools.GetUUID4(),
		conn:  conn,
		send:  make(chan []byte, 256),
		mu:    tools.NewTimedMutex("ClientMutex", 5),
		rooms: make(map[string]struct{}),
	}
}

func (c *Client) joinRoom(room string) {
	c.mu.Lock("joinRoom")
	c.rooms[room] = struct{}{}
	c.mu.Unlock()
	lgr.Info("Client", c.uid, "joined room:", room)
}

func (c *Client) leaveRoom(room string) {
	c.mu.Lock("leaveRoom")
	delete(c.rooms, room)
	c.mu.Unlock()
	lgr.Info("Client", c.uid, "left room:", room)
}

// ─────────────────────────────────────────────
//  Event handler registry
// ─────────────────────────────────────────────

var eventHandlers = map[string]func(*Client, Message){}

// ─────────────────────────────────────────────
//  WebSocket upgrader
// ─────────────────────────────────────────────

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: restrict in production
	},
}

// ─────────────────────────────────────────────
//  Connection handler
// ─────────────────────────────────────────────

func HandleWebSocketConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		lgr.Error("WebSocket upgrade error:", err)
		return
	}

	client := NewClient(conn)
	client.mu.Monitor(3 * time.Second)
	MYIO.addClient(client)

	go client.readPump()
	go client.writePump()

	lgr.Info("Client connected: %s", client.uid)
}

// ─────────────────────────────────────────────
//  readPump
// ─────────────────────────────────────────────

func (c *Client) readPump() {
	defer func() {
		MYIO.removeClient(c)
		c.conn.Close()
		lgr.Info("readPump exited for client: %s", c.uid)
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// Reset deadline every time we get a pong back from the client.
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure,
			) {
				log.Printf("Unexpected WebSocket close for %s: %v", c.uid, err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			lgr.Error("Invalid message format:", err)
			continue
		}

		// Built-in room management events
		switch msg.Event {
		case "join":
			if room, ok := msg.Data.(string); ok {
				c.joinRoom(room)
			}
		case "leave":
			if room, ok := msg.Data.(string); ok {
				c.leaveRoom(room)
			}
		case "join:session":
			// data: { sessionId: "abc" }
			if m, ok := msg.Data.(map[string]interface{}); ok {
				if sid, ok := m["sessionId"].(string); ok && sid != "" {
					c.joinRoom("session:" + sid)
				}
			}
		case "join:user":
			// data: { userUid: "xyz" }
			if m, ok := msg.Data.(map[string]interface{}); ok {
				if uid, ok := m["userUid"].(string); ok && uid != "" {
					c.joinRoom("user:" + uid)
				}
			}
		default:
			if handler, ok := eventHandlers[msg.Event]; ok {
				handler(c, msg)
			}
		}
	}
}

// ─────────────────────────────────────────────
//  writePump  (with ping ticker for health check)
// ─────────────────────────────────────────────

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		//lgr.Info("writePump exited for client:", c.uid)
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				lgr.Error("Write error for client", c.uid, ":", err)
				return // exit — readPump will clean up via removeClient
			}

		case <-ticker.C:
			// Send a ping; if the client doesn't pong within pongWait the
			// read deadline fires and readPump exits, which closes the conn.
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				lgr.Error("Ping failed for client: %s : %s", c.uid, err)
				return
			}
		}
	}
}

// ─────────────────────────────────────────────
//  emit helpers
// ─────────────────────────────────────────────

func (c *Client) emit(event string, data interface{}, ackId *int) {
	msg := Message{Event: event, Data: data, AckId: ackId}
	raw, err := json.Marshal(msg)
	if err != nil {
		lgr.ErrorStack(err.Error())
		return
	}
	// Non-blocking send; drop if the buffer is full (slow consumer).
	select {
	case c.send <- raw:
	default:
		lgr.Error("send buffer full, dropping message for client:", c.uid)
	}
}

func (c *Client) emitSync(event string, data interface{}, ackId *int) {
	msg := Message{Event: event, Data: data, AckId: ackId}
	raw, err := json.Marshal(msg)
	if err != nil {
		lgr.ErrorStack(err.Error())
		return
	}
	c.mu.Lock("emitSync")
	defer c.mu.Unlock()
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := c.conn.WriteMessage(websocket.TextMessage, raw); err != nil {
		lgr.Error("emitSync write error for client", c.uid, ":", err)
		c.conn.Close()
		MYIO.removeClient(c)
	}
}

// handleAcknowledgment replies to a client-initiated ack request.
func (c *Client) handleAcknowledgment(ackId int, ackData interface{}) {
	msg := Message{AckId: &ackId, AckData: ackData}
	raw, _ := json.Marshal(msg)
	c.send <- raw
}

// ─────────────────────────────────────────────
//  Server bootstrap
// ─────────────────────────────────────────────

func RunMyioServer(r *gin.Engine) {
	eventHandlers["message"] = func(c *Client, msg Message) {
		lgr.Info("Received message from", c.uid, ":", msg.Data)
		if msg.AckId != nil {
			c.handleAcknowledgment(*msg.AckId, "Server received your message!")
		}
	}

	r.GET("/my.io", HandleWebSocketConnection)
	lgr.Ok("My.io registered")
}

var MYIO = NewMyio()
