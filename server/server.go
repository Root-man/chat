package server

import (
	"errors"
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

		if len(s.conns) > 1 {
			msg := &packets.Message{From: "CHAT", Payload: fmt.Sprintf("User %s has joined the chat!", *username), Timestamp: time.Now()}
			// presence := &packets.Presence{Username: *username, Status: true}

			var to []string

			for u := range maps.Keys(s.conns) {
				to = append(to, u)
			}

			go s.multicast(msg, to)
			// go s.multicast(presence, to)
		}
		go s.handleConnection(*username)
	}
}

func (s *Server) handshake(conn net.Conn) (*string, error) {
	handshake := &packets.Handshake{}
	if err := handshake.Receive(conn); err != nil {
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
		err := msg.Receive(conn)
		if err != nil && err != io.EOF {
			log.Printf("Failed to deserialize message from %s: %s", username, err)
			continue
		} else if err == io.EOF {
			log.Printf("User %s connection closed by client", username)
			s.removeConnection(username)
			return
		}

		if len(s.conns) == 1 {
			continue
		}
		var to []string

		for u := range maps.Keys(s.conns) {
			if u != username {
				to = append(to, u)
			}
		}

		log.Printf("To: %v", to)

		s.multicast(msg, to)
	}
}

func (s *Server) removeConnection(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := &packets.Message{From: "CHAT", Payload: fmt.Sprintf("User %s has left the chat.", username), Timestamp: time.Now()}
	// presence := &packets.Presence{Username: username, Status: true}

	delete(s.conns, username)
	var to []string

	for u := range maps.Keys(s.conns) {
		to = append(to, u)
	}

	go s.multicast(msg, to)
	// go s.multicast(presence, to)
}

func (s *Server) multicast(p packets.Packet, to []string) error {
	for _, username := range to {
		if err := s.send(p, username); err != nil {
			return errors.Join(errors.New("failed to multicast packet"), err)
		}
	}

	return nil
}

func (s *Server) send(p packets.Packet, to string) error {
	log.Printf("Sending %s to %s", p, to)
	conn, ok := s.conns[to]
	if !ok {
		return fmt.Errorf("no connection found for %s", to)
	}

	log.Printf("Sending %s to %s", p, to)

	_, err := conn.Write(p.Encode())
	if err != nil {
		return errors.Join(fmt.Errorf("failed to send packet %s to %s", p, to), err)
	}

	log.Printf("Packet %s was sent to %s", p, to)

	return nil
}
