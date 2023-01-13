package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

const (
	Datagram  = 516          //maximum supported datagram size
	BlockSize = Datagram - 4 // Datagram minus 4 bytes header
)

type OpCode uint16

const (
	OpRRQ OpCode = iota + 1
	opData
	opAck
	opErr
)

type ErrCode uint16

const (
	ErrUnknown ErrCode = iota
	ErrNotFound
	ErrAccessViolation
	ErrDiskFull
	ErrIllegalOp
	ErrUnknownID
	ErrFileExists
	ErrNoUser
)

type ReadReq struct {
	Filename string
	Mode     string
}

func (q ReadReq) MarshalBinary() ([]byte, error) {
	mode := "octet"

	if q.Mode != "" {
		mode = q.Mode
	}

	// operation code + filename + 0 byte + mode + 0 byte ==== Request packet

	cap := 2 + 2 + len(q.Filename) + 1 + len(mode) + 1
	b := new(bytes.Buffer)
	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, OpRRQ) // writing operation code

	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(q.Filename) // writting filename

	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0) // writting 0 byte

	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(mode) // writting mode

	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0) // writting 0 byte

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil

}

func (q *ReadReq) UnMarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var code OpCode

	err := binary.Read(r, binary.BigEndian, &code)

	if err != nil {
		return err
	}

	if code != OpRRQ {
		return errors.New("invalid RRQ")
	}

	q.Filename, err = r.ReadString(0) // read filename

	if err != nil {
		return errors.New("invalid RRQ ")
	}

	q.Filename = strings.TrimRight(q.Filename, "\x00") //  remove 0 byte

	if len(q.Filename) == 0 {
		return errors.New("invalid RRQ ")
	}

	q.Mode, err = r.ReadString(0) // read mode

	if err != nil {
		return errors.New("invalid RRQ ")
	}

	q.Mode = strings.TrimRight(q.Mode, "\x00") //  remove 0 byte

	if len(q.Mode) == 0 {
		return errors.New("invalid RRQ ")
	}

	actual := strings.ToLower(q.Mode)

	if actual != "octet" {
		return errors.New("only binary transfer is required")
	}

	return nil
}
