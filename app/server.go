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
		log.Fatalln("ERROR: failed to start TCP listener: ", err.Error())
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatalln("ERROR: failed to accept connection: ", err.Error())
		}

		buf := make([]byte, 1024)

		n, err := conn.Read(buf)

		if err != nil {
			log.Println("ERROR: failed to read data from connection: ", err.Error())
		}

		command := string(buf[:n])

		if command == "*1\r\n$4\r\nPING\r\n" {
			PONG := "+PONG\r\n"

			_, err = conn.Write([]byte(PONG))

			if err != nil {
				log.Println("ERROR: failed to write data to connection: ", err.Error())
			}

			err = conn.Close()

			if err != nil {
				log.Println("ERROR: failed to close connection: ", err.Error())
			}
		}
	}
}
