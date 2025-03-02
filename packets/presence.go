package packets

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Presence struct {
	Username string
	Status   bool
}

func (p *Presence) Encode() []byte {
	usernameLength := uint32(len(p.Username))
	packet := make([]byte, 5+usernameLength)
	binary.BigEndian.PutUint32(packet[0:4], usernameLength)
	copy(packet[4:], []byte(p.Username))

	var statusByte byte

	if p.Status {
		statusByte = 1
	}

	packet[len(packet)-1] = statusByte
	return packet
}

func (p *Presence) Receive(r io.Reader) error {
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

	// read the status byte
	statusByte := make([]byte, 1)

	io.ReadFull(r, statusByte)

	p.Username = username
	if statusByte[0] == 1 {
		p.Status = true
	}

	return nil
}

func (p *Presence) String() string {
	return fmt.Sprintf("Presence: user %s is %v", p.Username, p.Status)
}
