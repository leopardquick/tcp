package tlv

import "errors"

const (
	BinaryType uint8 = iota + 1
	StringType

	MaxPayload uint32 = 10 << 20 //10MB
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceed")
