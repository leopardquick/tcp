package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Ack uint16

func (a *Ack) MarshalBinary() ([]byte, error) {

	b := new(bytes.Buffer)

	//opreation code  + Block number
	cap := 2 + 2

	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, opAck) // write operation code

	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, a) // write block number

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil

}

func (a *Ack) UnMarshalBinary(p []byte) error {

	r := bytes.NewBuffer(p)
	var opcode OpCode

	err := binary.Read(r, binary.BigEndian, &opcode) // read operation code

	if err != nil {
		return err
	}

	if opcode != opAck {
		return errors.New("")
	}

	return binary.Read(r, binary.BigEndian, a) // read block number

}
