package packets

import (
	"io"
)

type Packet interface {
	Encode() []byte
	Receive(io.Reader) error
	String() string
}
