package tlv

import (
	"encoding/binary"
	"io"
)

type String string

func (m String) String() string {
	return string(m)
}

func (m String) Bytes() []byte {
	return []byte(m)
}

func (m String) WriteTo(w io.Writer) (int64, error) {

	var n int64 = 0
	err := binary.Write(w, binary.BigEndian, StringType) // 1 byte

	if err != nil {
		return n, err
	}
	n += 1
	err = binary.Write(w, binary.BigEndian, int32(len(m))) // 4 - byte
	if err != nil {
		return n, err
	}

	n += 4
	o, err := w.Write([]byte(m)) // payload
	return int64(o) + n, err

}

func (m *String) ReadFrom(r io.Reader) (int64, error) {

	var typ uint8

	err := binary.Read(r, binary.BigEndian, &typ) // 1 -byte

	if err != nil {
		return 0, err
	}

	var n int64 = 1

	var size int32

	err = binary.Read(r, binary.BigEndian, &size) // 4 byte

	if err != nil {
		return n, err
	}
	n += 4

	buf := make([]byte, size)

	o, err := r.Read(buf) // payload

	if err != nil {
		return n, err
	}

	*m = String(buf)

	return int64(o) + n, err

}
