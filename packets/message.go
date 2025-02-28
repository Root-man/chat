package packets

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type Message struct {
	From      string
	Payload   string
	Timestamp time.Time
}

func (m *Message) String() string {
	return fmt.Sprintf("Message: From %s at %s", m.From, m.Timestamp)
}

func (m *Message) Encode() []byte {
	fromLength := uint32(len(m.From))
	messageLength := uint32(len(m.Payload))
	timestamp := m.Timestamp.Unix()

	packet := make([]byte, 16+fromLength+messageLength)

	binary.BigEndian.PutUint32(packet[0:4], fromLength)
	copy(packet[4:4+fromLength], []byte(m.From))
	binary.BigEndian.PutUint32(packet[4+fromLength:8+fromLength], messageLength)
	copy(packet[8+fromLength:8+fromLength+messageLength], []byte(m.Payload))
	binary.BigEndian.PutUint64(packet[8+fromLength+messageLength:], uint64(timestamp))

	return packet
}

func (m *Message) Decode(r io.Reader) error {
	// Read the first 4 bytes to get the length of the username
	lengthBytes := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBytes); err != nil {
		return err
	}

	usernameLength := binary.BigEndian.Uint32(lengthBytes)

	// Read the username based on the length
	usernameBytes := make([]byte, usernameLength)
	if _, err := io.ReadFull(r, usernameBytes); err != nil {
		return err
	}
	username := string(usernameBytes)

	// Read the next 4 bytes to get the length of the message
	if _, err := io.ReadFull(r, lengthBytes); err != nil {
		return err
	}
	messageLength := binary.BigEndian.Uint32(lengthBytes)

	// Read the message based on the length
	messageBytes := make([]byte, messageLength)
	if _, err := io.ReadFull(r, messageBytes); err != nil {
		return err
	}

	// Read the next 8 bytes to get the timestamp
	timestampBytes := make([]byte, 8)
	if _, err := io.ReadFull(r, timestampBytes); err != nil {
		return err
	}
	timestamp := time.Unix(int64(binary.BigEndian.Uint64(timestampBytes)), 0)

	m.From = username
	m.Payload = string(messageBytes)
	m.Timestamp = timestamp

	return nil
}
