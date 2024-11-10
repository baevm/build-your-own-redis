package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const ADDR = "0.0.0.0:6379"

func main() {
	log.Println("Starting TCP listener on:", ADDR)
	listener, err := net.Listen("tcp", ADDR)

	if err != nil {
		log.Fatalln("ERROR: failed to start TCP listener: ", err.Error())
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatalln("ERROR: failed to accept connection: ", err.Error())
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		command := scanner.Text()

		commandWords := strings.Split(command, "\r\n")
		fmt.Println(commandWords)

		for _, word := range commandWords {
			if strings.ToLower(word) == "ping" {
				PONG := "+PONG\r\n"

				_, err := conn.Write([]byte(PONG))

				if err != nil {
					log.Println("ERROR: failed to write data to connection: ", err.Error())
				}
			}
		}
	}
}
