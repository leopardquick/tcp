package tlv

import (
	"encoding/binary"
	"errors"
	"io"
)

type Binary []byte

func (m Binary) Bytes() []byte {
	return m
}

func (m Binary) String() string {
	return string(m)
}

func (m Binary) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, BinaryType) // 1-byte type
	if err != nil {
		return 0, err
	}

	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m)))
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write(m)

	return int64(o) + n, err

}

func (m *Binary) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8

	err := binary.Read(r, binary.BigEndian, &typ) // 1-byte

	if err != nil {
		return 0, err
	}

	var n int64 = 1

	if typ != BinaryType {
		return n, errors.New("not Binary")
	}

	var size uint32

	err = binary.Read(r, binary.BigEndian, &size) // 4 - byte
	if err != nil {
		return n, err
	}

	n += 4
	if size > MaxPayload {
		return n, ErrMaxPayloadSize
	}

	*m = make([]byte, size)

	o, err := r.Read(*m)
	return int64(o) + n, err
}
