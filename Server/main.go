package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	log := log.New(os.Stdin, "message ", 1)

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

				if err != nil {
					log.Print("error")
				}
				defer func() {
					con.Close()
				}()
				b := make([]byte, 1024)

				for {
					n, err := con.Read(b)

					if err != nil {
						if err == io.EOF {
							return
						} else {
							log.Fatal(err.Error())
						}
					}

					log.Printf("the message that have been read : %v", string(b[:n]))
				}
			}(conn)
		}
	}()

	<-done
	listener.Close()
}
