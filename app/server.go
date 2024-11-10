package main

import (
	"log"
	"net"
)

const ADDR = "0.0.0.0:6379"

func main() {
	log.Println("Starting TCP listener on:", ADDR)
	listener, err := net.Listen("tcp", ADDR)

	if err != nil {
		log.Fatalln("failed to start TCP listener: ", err.Error())
	}

	_, err = listener.Accept()

	if err != nil {
		log.Fatalln("failed to accept connection: ", err.Error())
	}
}
