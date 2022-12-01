package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {

	log := log.New(os.Stdin, "proxy", 0)

	done := make(chan struct{})

	server, err := net.Listen("tcp", "127.0.0.1:")

	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {

		for {
			conn, err := server.Accept()

			if err != nil {
				return
			}

			go func(con net.Conn) {

				buf := make([]byte, 1024)

				for {
					n, err := con.Read(buf)
					if err != nil {
						if err == io.EOF {
							log.Print(err)
						}
						return
					}

					switch msg := string(buf[:n]); msg {
					case "ping":
						_, err = con.Write([]byte("pong"))
					default:
						_, err = con.Write([]byte(string(buf[:n])))
					}

					if err != nil {
						if err == io.EOF {
							log.Print(err.Error())
						}
						return
					}

				}
			}(conn)
		}

	}()

	proxyServer, err := net.Listen("tcp", "127.0.0.1:")

	if err != nil {
		log.Fatal(err)
	}

	go func() {

		for {
			from, err := proxyServer.Accept()
			if err != nil {
				break
			}

			go func(from net.Conn) {
				defer from.Close()

				to, err := net.Dial("tcp", server.Addr().String())

				if err != nil {
					return
				}

				err = Proxy(from, to)

				if err != nil {
					return
				}

			}(from)

		}
	}()

	client, err := net.Dial("tcp", proxyServer.Addr().String())

	if err != nil {
		log.Fatal(err.Error())
	}

	defer client.Close()

	msg := []struct{ message, reply string }{{"ping", "pong"},
		{"pong", "pong"},
		{"echo", "echo"},
		{"ping", "pong"}}

	for i, m := range msg {
		_, err = client.Write([]byte(m.message))

		if err != nil {
			log.Fatal(err.Error())
		}

		b := make([]byte, 1024)

		nc, err := client.Read(b)

		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf(" %v  -> proxy ->  %v", msg[i].message, string(b[:nc]))
	}

	<-done

}

func ProxyConn(source, destinastion string) error {
	sourceConnection, err := net.Dial("tcp", source)

	if err != nil {
		return err
	}

	defer sourceConnection.Close()

	destinastionConnection, err := net.Dial("tcp", destinastion)

	if err != nil {
		return err
	}

	defer destinastionConnection.Close()

	go func() { _, _ = io.Copy(sourceConnection, destinastionConnection) }()

	//send message from source to destinatio

	_, err = io.Copy(destinastionConnection, sourceConnection)

	if err != nil {
		return err
	}

	return err

}

func Proxy(from io.Reader, to io.Writer) error {
	fromReader, isfromReader := from.(io.Writer)
	toWritter, isfromWritter := to.(io.Reader)

	if isfromReader && isfromWritter {
		go func() {
			_, _ = io.Copy(fromReader, toWritter)
		}()
	}

	_, err := io.Copy(to, from)

	return err
}
