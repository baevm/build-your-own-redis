package main

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/cache"
)

// Redis commands
const (
	ECHO = "echo"
	PING = "ping"

	SET = "set"
	GET = "get"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

func (w *Writer) Write(value Value, cache *cache.Cache) bool {
	i := 0

	// todo: remove arr loop?
	for i < len(value.arrayVal) {
		command := value.arrayVal[i]

		switch command.bulkStrVal {
		case PING:
			n, err := w.writePing()

			if err != nil {
				log.Println("ERROR: failed to PONG: ", err.Error())
			}

			i += n
			return true

		case ECHO:
			n, err := w.writeEcho(value, i)

			if err != nil {
				log.Println("ERROR: failed to write data to connection: ", err.Error())
			}

			i += n
			return true

		case SET:
			key := ""

			// read key
			if i+1 <= len(value.arrayVal) {
				i += 1
				key = value.arrayVal[i].bulkStrVal
			}

			val := ""

			// read value
			if i+1 <= len(value.arrayVal) {
				i += 1
				val = value.arrayVal[i].bulkStrVal
			}

			// if PX argument with time provided
			// set cache item with expire time
			// else set item with no expiration
			if i+2 <= len(value.arrayVal) {
				// read PX argument
				i += 2
				expireTimeStr := value.arrayVal[i].bulkStrVal

				expireTimeMs, err := strconv.Atoi(expireTimeStr)

				if err != nil {
					log.Println("Failed to convert expire time to int: ", err.Error())
					return false
				}

				cache.SetWithExpiration(key, val, expireTimeMs)
			} else {
				cache.Set(key, val)
			}

			OK := "+OK\r\n"

			_, err := w.writer.Write([]byte(OK))

			if err != nil {
				log.Println("ERROR: failed to OK after SET: ", err.Error())
				return false
			}

			return true

		case GET:
			key := ""

			if i+1 <= len(value.arrayVal) {
				i += 1
				key = value.arrayVal[i].bulkStrVal
			}

			val, isFound := cache.Get(key)

			if !isFound {
				NOTFOUND := "$-1\r\n"

				_, err := w.writer.Write([]byte(NOTFOUND))

				if err != nil {
					log.Println("ERROR: failed to NOTFOUND in GET: ", err.Error())
					return false
				}

				return true
			}

			encodedStr := fmt.Sprintf("$%v\r\n%s\r\n", len(val), val)

			_, err := w.writer.Write([]byte(encodedStr))

			if err != nil {
				log.Println("ERROR: failed to response to GET: ", err.Error())
				return false
			}

			i += 1
			return true

		default:
			return false
		}
	}

	return false
}

// Answer to PING command
func (w *Writer) writePing() (int, error) {
	PONG := "+PONG\r\n"

	_, err := w.writer.Write([]byte(PONG))

	return 1, err
}

// Answer to ECHO command
func (w *Writer) writeEcho(value Value, i int) (int, error) {
	currI := i

	echoMsg := ""

	// Read message of ECHO command
	// it should be next in array
	if currI+1 <= len(value.arrayVal)-1 {
		currI += 1
		echoMsg = value.arrayVal[currI].bulkStrVal
	}

	encodedStr := fmt.Sprintf("$%v\r\n%s\r\n", len(echoMsg), echoMsg)

	_, err := w.writer.Write([]byte(encodedStr))

	if currI-i == 0 {
		return 1, err
	} else {
		return currI - i, err
	}
}
