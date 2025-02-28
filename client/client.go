package client

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/root-man/chat/packets"
)

type Client struct {
	name        string
	conn        net.Conn
	usersOnline []string
}

func New(name string) *Client {
	return &Client{name: name}
}

func (c *Client) Connect(serverHost string, serverPort int) (<-chan packets.Message, error) {
	log.Printf("Client %s connecting to chat server %s:%d", c.name, serverHost, serverPort)
	tcpServer, err := net.ResolveTCPAddr("tcp", serverHost+":"+strconv.Itoa(serverPort))
	if err != nil {
		log.Printf("ResolveTCPAddr failed: %s", err)
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		log.Printf("Dial failed: %s", err)
		return nil, err
	}

	c.conn = conn

	log.Printf("Connection successful, initiating handshake...")

	if err := c.handshake(); err != nil {
		log.Printf("error on handshake: %s", err)
		return nil, err
	}

	log.Printf("Handshake completed")

	msgChan := make(chan packets.Message)

	go c.listen(msgChan)

	return msgChan, nil
}

func (c *Client) handshake() error {
	h := packets.Handshake{Username: c.name}
	_, err := c.conn.Write(h.Encode())
	if err != nil {
		return fmt.Errorf("handshake failed: %s", err)
	}

	resp := &packets.HandshakeResponse{}

	if err = resp.Decode(c.conn); err != nil {
		return fmt.Errorf("handshake failed: %s", err)
	}

	c.usersOnline = resp.OnlineUsers

	return nil
}

func (c *Client) listen(msgChan chan<- packets.Message) {
	defer c.conn.Close()

	msg := &packets.Message{}

	for {
		if err := msg.Decode(c.conn); err != nil {
			fmt.Printf("Failed to deserialize message: %s", err)
			continue
		}

		msgChan <- *msg
	}
}

func (c *Client) Send(message string) (*packets.Message, error) {
	msg := &packets.Message{From: c.name, Payload: message, Timestamp: time.Now()}
	_, err := c.conn.Write(msg.Encode())
	if err != nil {
		return nil, err
	}

	return msg, nil
}
