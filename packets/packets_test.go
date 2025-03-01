package packets

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPresence_Encode(t *testing.T) {
	p := Presence{
		Username: "testuser",
		Status:   true,
	}

	expected := []byte{
		0, 0, 0, 8, // Length of the username (8 bytes)
		't', 'e', 's', 't', 'u', 's', 'e', 'r', // Username
		1, // Status (true)
	}

	encoded := p.Encode()

	if !bytes.Equal(encoded, expected) {
		t.Errorf("Encode() = %v, want %v", encoded, expected)
	}
}

func TestPresence_Decode(t *testing.T) {
	data := []byte{
		0, 0, 0, 8, // Length of the username (8 bytes)
		't', 'e', 's', 't', 'u', 's', 'e', 'r', // Username
		1, // Status (true)
	}

	r := bytes.NewReader(data)
	var p Presence
	if err := p.Decode(r); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	expected := Presence{
		Username: "testuser",
		Status:   true,
	}

	if p != expected {
		t.Errorf("Decode() = %v, want %v", p, expected)
	}
}

func TestHandshake_Encode(t *testing.T) {
	h := Handshake{
		Username: "testuser",
	}

	expected := []byte{
		0, 0, 0, 8, // Length of the username (8 bytes)
		't', 'e', 's', 't', 'u', 's', 'e', 'r', // Username
	}

	encoded := h.Encode()

	if !bytes.Equal(encoded, expected) {
		t.Errorf("Encode() = %v, want %v", encoded, expected)
	}
}

func TestHandshake_Decode(t *testing.T) {
	data := []byte{
		0, 0, 0, 8, // Length of the username (8 bytes)
		't', 'e', 's', 't', 'u', 's', 'e', 'r', // Username
	}

	r := bytes.NewReader(data)
	var h Handshake
	if err := h.Decode(r); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	expected := Handshake{
		Username: "testuser",
	}

	if h != expected {
		t.Errorf("Decode() = %v, want %v", h, expected)
	}
}

func TestHandshakeResponse_Encode(t *testing.T) {
	hr := HandshakeResponse{
		OnlineUsers: []string{"user1", "user2"},
	}

	expected := []byte{
		0, 0, 0, 2, // Number of online users (2)
		0, 0, 0, 5, // Length of the first username (5 bytes)
		'u', 's', 'e', 'r', '1', // First username
		0, 0, 0, 5, // Length of the second username (5 bytes)
		'u', 's', 'e', 'r', '2', // Second username
	}

	encoded := hr.Encode()

	if !bytes.Equal(encoded, expected) {
		t.Errorf("Encode() = %v, want %v", encoded, expected)
	}
}

func TestHandshakeResponse_Decode(t *testing.T) {
	data := []byte{
		0, 0, 0, 2, // Number of online users (2)
		0, 0, 0, 5, // Length of the first username (5 bytes)
		'u', 's', 'e', 'r', '1', // First username
		0, 0, 0, 5, // Length of the second username (5 bytes)
		'u', 's', 'e', 'r', '2', // Second username
	}

	r := bytes.NewReader(data)
	var hr HandshakeResponse
	if err := hr.Decode(r); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	expected := HandshakeResponse{
		OnlineUsers: []string{"user1", "user2"},
	}

	if !assert.Equal(t, hr, expected) {
		t.Errorf("Expected %v, got %v", expected, hr)
	}
}

func TestMessage_Encode(t *testing.T) {
	m := Message{
		From:      "testuser",
		Payload:   "Hello, world!",
		Timestamp: time.Unix(256, 0), // Example timestamp
	}

	expected := []byte{
		0, 0, 0, 8, // Length of the username (8 bytes)
		't', 'e', 's', 't', 'u', 's', 'e', 'r', // Username
		0, 0, 0, 13, // Length of the message (13 bytes)
		'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', '!', // Message
		0, 0, 0, 0, 0, 0, 1, 0, // Timestamp (256)
	}

	encoded := m.Encode()

	if !bytes.Equal(encoded, expected) {
		t.Errorf("Encode() = %v, want %v", encoded, expected)
	}
}

func TestMessage_Decode(t *testing.T) {
	data := []byte{
		0, 0, 0, 8, // Length of the username (8 bytes)
		't', 'e', 's', 't', 'u', 's', 'e', 'r', // Username
		0, 0, 0, 13, // Length of the message (13 bytes)
		'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', '!', // Message
		0, 0, 0, 0, 0, 0, 1, 1, // Timestamp (257)
	}

	r := bytes.NewReader(data)
	var m Message
	if err := m.Decode(r); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	expected := Message{
		From:      "testuser",
		Payload:   "Hello, world!",
		Timestamp: time.Unix(257, 0), // Example timestamp
	}

	if m != expected {
		t.Errorf("Decode() = %v, want %v", m, expected)
	}
}
