package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const ADDR = "0.0.0.0:6379"

func main() {
	redisServer := NewRedisServer(ADDR)
	redisServer.Start()
}

type RedisServer struct {
	addr string
}

func NewRedisServer(addr string) *RedisServer {
	return &RedisServer{
		addr: addr,
	}
}

func (rs *RedisServer) Start() {
	log.Println("Starting Redis server on:", rs.addr)
	listener, err := net.Listen("tcp", rs.addr)

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

	resp := NewResp(conn)

	for {
		value, err := resp.Read()

		if err != nil {
			if err == io.EOF {
				continue
			} else {
				log.Println("fail to read RESP: ", err)
				return
			}
		}

		log.Printf("%+v\n", value)

		if value.dataType == "bulk" {

		} else {
			i := 0

			for i < len(value.arrayVal) {
				command := value.arrayVal[i]

				switch command.bulkStrVal {
				case "ping":
					PONG := "+PONG\r\n"

					_, err := conn.Write([]byte(PONG))

					if err != nil {
						log.Println("ERROR: failed to write data to connection: ", err.Error())
					}

					i += 1

				case "echo":
					echoMsg := ""

					if i+1 <= len(value.arrayVal)-1 {
						i += 1
						echoMsg = value.arrayVal[i].bulkStrVal
					}

					encodedStr := fmt.Sprintf("$%v\r\n%s\r\n", len(echoMsg), echoMsg)

					_, err := conn.Write([]byte(encodedStr))

					if err != nil {
						log.Println("ERROR: failed to write data to connection: ", err.Error())
					}

					i += 1

				default:
					return
				}
			}
		}
	}

}
