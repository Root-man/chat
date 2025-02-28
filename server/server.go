package server

import (
	"fmt"
	"io"
	"log"
	"maps"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/root-man/chat/packets"
)

type Server struct {
	listener net.Listener
	conns    map[string]net.Conn
	mu       sync.Mutex
}

func New(portNumber int) (*Server, error) {
	PORT := ":" + strconv.Itoa(portNumber)
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Printf("Started chat server listening on port %d", portNumber)

	return &Server{listener: l, mu: sync.Mutex{}, conns: make(map[string]net.Conn)}, nil
}

func (s *Server) Run() error {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err)
			return err
		}

		log.Printf("Got incoming connection from %s, initiating handshake...", c.RemoteAddr())
		username, err := s.handshake(c)
		if err != nil {
			log.Printf("Handshake with %s failed: %s", *username, err)
			c.Close()
			continue
		}

		msg := &packets.Message{From: "CHAT", Payload: fmt.Sprintf("User %s has joined the chat!", *username), Timestamp: time.Now()}

		to := make([]string, len(s.conns))

		for u := range maps.Keys(s.conns) {
			to = append(to, u)
		}

		go s.broadcast(msg, to)
		go s.handleConnection(*username)
	}
}

func (s *Server) handshake(conn net.Conn) (*string, error) {
	handshake := &packets.Handshake{}
	if err := handshake.Decode(conn); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.conns[handshake.Username]; ok {
		return nil, fmt.Errorf("username %s is already in use", handshake.Username)
	}

	onlineUsers := make([]string, len(s.conns))
	for i := range maps.Keys(s.conns) {
		onlineUsers = append(onlineUsers, i)
	}

	handshakeResponse := packets.HandshakeResponse{OnlineUsers: onlineUsers}

	_, err := conn.Write(handshakeResponse.Encode())
	if err != nil {
		return nil, err
	}

	s.conns[handshake.Username] = conn
	log.Printf("Handshake successful with username: %s", handshake.Username)
	return &handshake.Username, nil
}

func (s *Server) handleConnection(username string) {
	conn := s.conns[username]
	defer conn.Close()

	msg := &packets.Message{}

	for {
		err := msg.Decode(conn)
		if err != nil && err != io.EOF {
			log.Printf("Failed to deserialize message from %s: %s", username, err)
			continue
		} else if err == io.EOF {
			log.Printf("User %s connection closed by client", username)
			s.removeConnection(username)
			return
		}

		to := make([]string, len(s.conns)-1)

		for u := range maps.Keys(s.conns) {
			if u != username {
				to = append(to, u)
			}
		}

		s.broadcast(msg, to)
	}
}

func (s *Server) removeConnection(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := &packets.Message{From: "CHAT", Payload: fmt.Sprintf("User %s has left the chat.", username), Timestamp: time.Now()}

	delete(s.conns, username)
	to := make([]string, len(s.conns))

	for u := range maps.Keys(s.conns) {
		to = append(to, u)
	}

	go s.broadcast(msg, to)
}

func (s *Server) broadcast(p packets.Packet, to []string) error {
	for _, username := range to {
		conn := s.conns[username]
		_, err := conn.Write(p.Encode())
		if err != nil {
			return fmt.Errorf("failed to broadcast the message to %s: %s", username, err)
		}

		log.Printf("Broadcasted %s to %s", p, username)
	}

	return nil
}
