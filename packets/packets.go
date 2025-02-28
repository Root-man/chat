package packets

import (
	"io"
)

type Packet interface {
	Encode() []byte
	Decode(io.Reader) error
	String() string
}
