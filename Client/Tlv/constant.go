package tlv

import "errors"

const (
	BinaryType = iota + 1
	StringType

	MaxPayload uint32 = 10 << 20 //10MB
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceed")
