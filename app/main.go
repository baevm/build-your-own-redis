package main

import "github.com/codecrafters-io/redis-starter-go/app/cache"

const ADDR = "0.0.0.0:6379"

func main() {
	cache := cache.CreateCache()

	redisServer := NewRedisServer(ADDR, cache)
	redisServer.Start()
}
