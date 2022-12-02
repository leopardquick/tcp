package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {

	file, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	monitor := Monitor{log.New(file, "monitor ", 0)}

	done := make(chan struct{})
	server, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		monitor.Fatal(err)
	}

	go func() {
		for {
			conn, err := server.Accept()

			if err != nil {
				monitor.Fatal(err.Error())
			}

			go func(con net.Conn) {
				defer func() {
					done <- struct{}{}
				}()
				defer con.Close()
				buffer := make([]byte, 1024)
				for {
					r := io.TeeReader(con, &monitor)

					n, err := r.Read(buffer)

					if err != nil {
						if err == io.EOF {
							monitor.Print(err)
						}

					}

					writer := io.MultiWriter(con, &monitor)

					_, err = writer.Write(buffer[:n])
					if err != nil {
						if err == io.EOF {
							monitor.Print(err.Error())
						}
						return
					}
				}
			}(conn)

		}
	}()

	client, err := net.Dial("tcp", server.Addr().String())

	if err != nil {
		monitor.Fatal(err)
	}

	go func(c net.Conn) {
		defer func() {
			c.Close()
			done <- struct{}{}
		}()

		buffer := make([]byte, 1024)

		_, err := c.Write([]byte("mambo"))

		if err != nil {
			if err == io.EOF {

			}
			return
		}

		_, err = c.Read(buffer)

		if err != nil {
			if err == io.EOF {
				monitor.Print(err.Error())
			}
			return
		}

	}(client)

	<-done
}
