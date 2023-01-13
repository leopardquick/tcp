package main

import (
	"io/ioutil"
	"log"
	"net"
)

func main() {

	payload, err := ioutil.ReadFile("test.txt")

	if err != nil {
		log.Fatal(err.Error())
	}

	server, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("listerning at port : %v", server.Addr().String())

	conn, err := server.Accept()

	if err != nil {
		log.Fatal(err)
	}

	conn.Write(payload)

	conn.Close()
	server.Close()

}
