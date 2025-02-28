package packets

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Handshake struct {
	Username string
}

func (h *Handshake) String() string {
	return fmt.Sprintf("Handshake: user %s", h.Username)
}

func (h *Handshake) Encode() []byte {
	usernameLength := uint32(len(h.Username))
	packet := make([]byte, 4+usernameLength)
	binary.BigEndian.PutUint32(packet[0:4], usernameLength)
	copy(packet[4:], []byte(h.Username))
	return packet
}

func (h *Handshake) Decode(r io.Reader) error {
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

	h.Username = username

	return nil
}
