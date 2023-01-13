package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

type Err struct {
	Error   ErrCode
	Message string
}

func (e *Err) MarshalBinary() ([]byte, error) {

	b := new(bytes.Buffer)

	// operation code + error code+ message + 0 bytes

	cap := 2 + 2 + len(e.Message) + 1

	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, opErr) // write operation code

	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, e.Error) // write error code

	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(e.Message) // write message

	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0) // write 0 byte

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil

}

func (e *Err) UnMarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var opcode OpCode

	err := binary.Read(r, binary.BigEndian, &opcode) // read operation cod

	if err != nil {
		return err
	}

	if opcode != opErr {
		return errors.New("invalid Error")
	}

	err = binary.Read(r, binary.BigEndian, &e.Error) // read error code

	if err != nil {
		return err
	}

	e.Message, err = r.ReadString(0)
	if err != nil {
		return err
	}

	e.Message = strings.TrimRight(e.Message, "\x00") // remove 0 - byte

	return nil
}
