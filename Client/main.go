package main

import (
	"context"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	. "github.com/leopardquick/tcp/client/Tlv"
)

func main() {

	log := log.New(os.Stdin, "message ", 1)

	dl := time.Now().Add(5 * time.Second)
	contx, cancel := context.WithDeadline(context.Background(), dl)

	defer cancel()

	var d net.Dialer
	d.Control = func(_, address string, c syscall.RawConn) error {
		return nil
	}

	conn, err := d.DialContext(contx, "tcp", "127.0.0.1:8080")

	if err != nil {
		nError, ok := err.(net.Error)
		if !ok {
			log.Fatal("Not temporary")
		} else {
			if !nError.Timeout() {
				log.Fatal("Error not timeout")
			}
		}
	}

	if contx.Err() == context.DeadlineExceeded {
		log.Fatal("context deadline exceed")
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	b1 := new(Binary)

	if err != nil {
		log.Fatal(err)
	}

	for {
		n, err := b1.ReadFrom(conn)

		if err != nil {
			log.Fatal(err.Error())
			break
		}
		log.Printf("value : %v b1 : %v", n, b1)
	}

	conn.Close()

}

// func Decoder(r io.Reader) (Payload, error) {

// 	var typ uint8
// 	err := binary.Read(r, binary.BigEndian, &typ)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var payload Payload

// 	switch typ {
// 	case BinaryType:
// 		payload = new(Binary)
// 	case StringType:
// 		payload = new(String)
// 	}

// 	_, err = payload.ReadFrom(
// 		io.MultiReader(bytes.NewReader([]byte{typ}), r))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return payload, nil
// }
