package main

import (
	"log"
	"net"
	"os"

	. "github.com/leopardquick/tcp/server/Tlv"
)

func main() {
	log := log.New(os.Stdin, "message ", 1)

	b1 := Binary("wheree to")

	payloads := []Payload{&b1}

	listener, err := net.Listen("tcp", "127.0.0.1:8080")

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	done := make(chan struct{})

	go func() {
		defer func() {
			done <- struct{}{}
		}()
		log.Printf("server listening at %v", listener.Addr())
		for {
			conn, err := listener.Accept()

			if err != nil {
				log.Fatal(err.Error())
				return
			}
			go func(con net.Conn) {

				for i := range payloads {
					n, e := payloads[i].WriteTo(con)
					if e != nil {
						log.Print(n)
						log.Fatal(e.Error())
					}
				}

			}(conn)
		}
	}()

	<-done
	listener.Close()
}
