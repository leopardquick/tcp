package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")

	if err != nil {
		log.Fatal(err.Error())
	}

	b := make([]byte, 1024)

	n, err := conn.Read(b)

	if err != nil {
		if err == io.EOF {

		} else {
			log.Fatal(err.Error())
		}
	}

	ioutil.WriteFile("test.txt", b[:n], 0666)
	fmt.Print(b[:n])
}
