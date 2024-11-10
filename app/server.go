package main

import (
	"io"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/cache"
)

type RedisServer struct {
	addr  string
	cache *cache.Cache
}

func NewRedisServer(addr string, cache *cache.Cache) *RedisServer {
	return &RedisServer{
		addr:  addr,
		cache: cache,
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

		go rs.handleConn(conn)
	}
}

func (rs *RedisServer) handleConn(conn net.Conn) {
	defer conn.Close()

	resp := NewResp(conn)

	for {
		value, err := resp.Read()

		if err != nil {

			if err == io.EOF {
				return
			} else {
				log.Println("failed to read RESP command: ", err)
				return
			}
		}

		log.Printf("%+v\n", value)

		writer := NewWriter(conn)
		isWrote := writer.Write(value, rs.cache)

		if !isWrote {
			return
		}
	}
}
