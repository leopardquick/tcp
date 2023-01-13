package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Data struct {
	Block   uint16
	Payload io.Reader
}

func (d *Data) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	b.Grow(Datagram)

	d.Block++

	err := binary.Write(b, binary.BigEndian, opData) //write opration head

	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, d.Block) //write block number

	if err != nil {
		return nil, err
	}

	_, err = io.CopyN(b, d.Payload, BlockSize)

	if err != nil && err != io.EOF {
		return nil, err
	}

	return b.Bytes(), nil
}

func (d *Data) UnMarshal(p []byte) error {
	if l := len(p); l < 4 || l > Datagram {
		return errors.New("Invalid Data")
	}

	var opcode OpCode

	err := binary.Read(bytes.NewBuffer(p[:2]), binary.BigEndian, &opcode)

	if err != nil || opcode != opData {
		return errors.New("invalid Data")
	}

	err = binary.Read(bytes.NewBuffer(p[2:4]), binary.BigEndian, &d.Block)

	if err != nil {
		return errors.New("inValid Data")
	}

	d.Payload = bytes.NewBuffer(p[4:])

	return nil
}
