package realtime

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const roomID = "investor-admin"

type Handler struct {
	hub *Hub
}

type Hub struct {
	mu            sync.Mutex
	clients       map[*Client]bool
	messages      []Event
	notifications []Event
}

type Client struct {
	conn net.Conn
	role string
	name string
	send chan Event
	hub  *Hub
}

type Event struct {
	Type       string `json:"type"`
	ID         string `json:"id,omitempty"`
	RoomID     string `json:"roomId,omitempty"`
	SenderID   string `json:"senderId,omitempty"`
	SenderName string `json:"senderName,omitempty"`
	SenderRole string `json:"senderRole,omitempty"`
	Body       string `json:"body,omitempty"`
	Title      string `json:"title,omitempty"`
	CreatedAt  string `json:"createdAt,omitempty"`
	Read       bool   `json:"read,omitempty"`
}

type inboundEvent struct {
	Type       string `json:"type"`
	RoomID     string `json:"roomId"`
	SenderID   string `json:"senderId"`
	SenderName string `json:"senderName"`
	SenderRole string `json:"senderRole"`
	Body       string `json:"body"`
}

func NewRealtimeHandler() *Handler {
	return &Handler{hub: &Hub{clients: map[*Client]bool{}}}
}

func (h *Handler) RegisterRealtimeHandler(router *gin.RouterGroup) {
	router.GET("/ws", h.handleWebSocket)
	router.GET("/messages", h.listMessages)
	router.POST("/messages", h.createMessage)
}

func (h *Handler) listMessages(c *gin.Context) {
	c.JSON(http.StatusOK, h.hub.snapshot())
}

func (h *Handler) createMessage(c *gin.Context) {
	var inbound inboundEvent
	if err := c.ShouldBindJSON(&inbound); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message, notification, ok := h.hub.buildMessage(inbound, inbound.SenderRole, inbound.SenderName)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message body is required"})
		return
	}
	h.hub.broadcast(message)
	h.hub.broadcast(notification)
	c.JSON(http.StatusCreated, gin.H{"message": message, "notification": notification})
}

func (h *Handler) handleWebSocket(c *gin.Context) {
	conn, rw, err := upgrade(c.Writer, c.Request)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	client := &Client{
		conn: conn,
		role: c.Query("role"),
		name: c.Query("name"),
		send: make(chan Event, 32),
		hub:  h.hub,
	}
	if client.name == "" {
		client.name = "Guest"
	}
	if client.role == "" {
		client.role = "guest"
	}

	h.hub.add(client)
	go client.writeLoop(rw.Writer)
	client.readLoop(rw.Reader)
}

func (h *Hub) add(client *Client) {
	h.mu.Lock()
	h.clients[client] = true
	snapshot := Event{Type: "snapshot", Body: mustJSON(h.snapshotLocked())}
	h.mu.Unlock()
	client.send <- snapshot
}

func (h *Hub) snapshot() map[string][]Event {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.snapshotLocked()
}

func (h *Hub) snapshotLocked() map[string][]Event {
	messages := append([]Event(nil), h.messages...)
	notifications := append([]Event(nil), h.notifications...)
	return map[string][]Event{"messages": messages, "notifications": notifications}
}

func (h *Hub) buildMessage(inbound inboundEvent, clientRole string, clientName string) (Event, Event, bool) {
	if strings.TrimSpace(inbound.Body) == "" {
		return Event{}, Event{}, false
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	senderRole := firstNonEmpty(inbound.SenderRole, clientRole)
	senderName := firstNonEmpty(inbound.SenderName, clientName)
	message := Event{
		Type:       "message",
		ID:         fmt.Sprintf("msg_%d", time.Now().UnixNano()),
		RoomID:     firstNonEmpty(inbound.RoomID, roomID),
		SenderID:   inbound.SenderID,
		SenderName: senderName,
		SenderRole: senderRole,
		Body:       strings.TrimSpace(inbound.Body),
		CreatedAt:  now,
	}
	notification := Event{
		Type:       "notification",
		ID:         fmt.Sprintf("ntf_%d", time.Now().UnixNano()),
		RoomID:     message.RoomID,
		SenderID:   message.SenderID,
		SenderName: senderName,
		SenderRole: senderRole,
		Title:      fmt.Sprintf("New message from %s", senderName),
		Body:       message.Body,
		CreatedAt:  now,
	}
	return message, notification, true
}

func (h *Hub) remove(client *Client) {
	h.mu.Lock()
	if h.clients[client] {
		delete(h.clients, client)
		close(client.send)
	}
	h.mu.Unlock()
	_ = client.conn.Close()
}

func (h *Hub) broadcast(event Event) {
	h.mu.Lock()
	switch event.Type {
	case "message":
		h.messages = append(h.messages, event)
		if len(h.messages) > 100 {
			h.messages = h.messages[len(h.messages)-100:]
		}
	case "notification":
		h.notifications = append([]Event{event}, h.notifications...)
		if len(h.notifications) > 50 {
			h.notifications = h.notifications[:50]
		}
	}
	for client := range h.clients {
		select {
		case client.send <- event:
		default:
		}
	}
	h.mu.Unlock()
}

func (c *Client) readLoop(reader *bufio.Reader) {
	defer c.hub.remove(c)
	for {
		payload, opcode, err := readFrame(reader)
		if err != nil {
			return
		}
		if opcode == 8 {
			return
		}
		if opcode != 1 {
			continue
		}
		var inbound inboundEvent
		if err := json.Unmarshal(payload, &inbound); err != nil {
			continue
		}
		message, notification, ok := c.hub.buildMessage(inbound, c.role, c.name)
		if !ok {
			continue
		}
		c.hub.broadcast(message)
		c.hub.broadcast(notification)
	}
}

func (c *Client) writeLoop(writer *bufio.Writer) {
	for event := range c.send {
		payload, err := json.Marshal(event)
		if err != nil {
			continue
		}
		if err := writeTextFrame(writer, payload); err != nil {
			return
		}
	}
}

func upgrade(w http.ResponseWriter, r *http.Request) (net.Conn, *bufio.ReadWriter, error) {
	if !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		return nil, nil, errors.New("missing websocket upgrade header")
	}
	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		return nil, nil, errors.New("missing websocket key")
	}
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("websocket hijacking is not supported")
	}
	conn, rw, err := hijacker.Hijack()
	if err != nil {
		return nil, nil, err
	}
	accept := websocketAccept(key)
	response := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + accept + "\r\n\r\n"
	if _, err := rw.WriteString(response); err != nil {
		_ = conn.Close()
		return nil, nil, err
	}
	if err := rw.Flush(); err != nil {
		_ = conn.Close()
		return nil, nil, err
	}
	return conn, rw, nil
}

func websocketAccept(key string) string {
	sum := sha1.Sum([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	return base64.StdEncoding.EncodeToString(sum[:])
}

func readFrame(reader *bufio.Reader) ([]byte, byte, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, 0, err
	}
	opcode := header[0] & 0x0f
	masked := header[1]&0x80 != 0
	length := uint64(header[1] & 0x7f)
	switch length {
	case 126:
		extended := make([]byte, 2)
		if _, err := io.ReadFull(reader, extended); err != nil {
			return nil, 0, err
		}
		length = uint64(binary.BigEndian.Uint16(extended))
	case 127:
		extended := make([]byte, 8)
		if _, err := io.ReadFull(reader, extended); err != nil {
			return nil, 0, err
		}
		length = binary.BigEndian.Uint64(extended)
	}
	var mask [4]byte
	if masked {
		if _, err := io.ReadFull(reader, mask[:]); err != nil {
			return nil, 0, err
		}
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return nil, 0, err
	}
	if masked {
		for i := range payload {
			payload[i] ^= mask[i%4]
		}
	}
	return payload, opcode, nil
}

func writeTextFrame(writer *bufio.Writer, payload []byte) error {
	header := []byte{0x81}
	length := len(payload)
	switch {
	case length < 126:
		header = append(header, byte(length))
	case length <= 65535:
		header = append(header, 126, byte(length>>8), byte(length))
	default:
		header = append(header, 127)
		size := make([]byte, 8)
		binary.BigEndian.PutUint64(size, uint64(length))
		header = append(header, size...)
	}
	if _, err := writer.Write(header); err != nil {
		return err
	}
	if _, err := writer.Write(payload); err != nil {
		return err
	}
	return writer.Flush()
}

func mustJSON(value any) string {
	payload, _ := json.Marshal(value)
	return string(payload)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
