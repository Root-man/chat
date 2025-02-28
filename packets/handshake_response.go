package packets

import (
	"encoding/binary"
	"fmt"
	"io"
)

type HandshakeResponse struct {
	OnlineUsers []string
}

func (hr *HandshakeResponse) String() string {
	return fmt.Sprintf("Handshake response: online users %s", hr.OnlineUsers)
}

func (hr *HandshakeResponse) Encode() []byte {
	numUsers := uint32(len(hr.OnlineUsers))
	packet := make([]byte, 4)

	binary.BigEndian.PutUint32(packet[0:4], numUsers)

	for _, user := range hr.OnlineUsers {
		usernameLength := uint32(len(user))
		userPacket := make([]byte, 4+usernameLength)
		binary.BigEndian.PutUint32(userPacket[0:4], usernameLength)
		copy(userPacket[4:], []byte(user))
		packet = append(packet, userPacket...)
	}

	return packet
}

func (hr *HandshakeResponse) Decode(r io.Reader) error {
	// Read the first 4 bytes to get the number of online users
	numUsersBytes := make([]byte, 4)
	if _, err := io.ReadFull(r, numUsersBytes); err != nil {
		return err
	}

	numUsers := binary.BigEndian.Uint32(numUsersBytes)
	onlineUsers := make([]string, numUsers)

	for i := uint32(0); i < numUsers; i++ {
		// Read the next 4 bytes to get the length of the username
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
		onlineUsers[i] = string(usernameBytes)
	}

	hr.OnlineUsers = onlineUsers

	return nil
}
