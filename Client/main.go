package main

import (
	"context"
	"log"
	"net"
	"os"
	"syscall"
	"time"
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

	conn.Write([]byte("server receive"))

	defer conn.Close()
}
