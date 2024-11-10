package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	// https://redis.io/docs/latest/develop/reference/protocol-spec/#arrays
	ARRAY = '*'

	// https://redis.io/docs/latest/develop/reference/protocol-spec/#bulk-strings
	BULK_STR = '$'

	// https://redis.io/docs/latest/develop/reference/protocol-spec/#simple-strings
	SIMPLE_STR = '+'
)

type Value struct {
	// 'array', 'bulk'
	dataType string

	// value of bulk string
	// in lower case
	bulkStrVal string

	// value of array
	arrayVal []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(conn net.Conn) *Resp {
	return &Resp{
		reader: bufio.NewReader(conn),
	}
}

func (resp *Resp) Read() (Value, error) {
	dataType, err := resp.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch dataType {
	case ARRAY:
		return resp.readArray()

	case BULK_STR:
		return resp.readBulk()

	default:
		log.Println("Unknown type: ", string(dataType))
		return Value{}, err
	}
}

// Reads RESP "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
func (resp *Resp) readArray() (Value, error) {
	val := Value{
		dataType: "array",
	}

	// size of array
	n, _, err := resp.readInteger()

	if err != nil {
		return val, err
	}

	val.arrayVal = make([]Value, n)

	for i := 0; i < n; i++ {
		newVal, err := resp.Read()

		if err != nil {
			return val, err
		}

		val.arrayVal[i] = newVal
	}

	return val, nil
}

// Reads RESP "$10\r\nhelloworld\r\n"
func (resp *Resp) readBulk() (Value, error) {
	val := Value{
		dataType: "bulk",
	}

	// Length of string
	n, _, err := resp.readInteger()

	if err != nil {
		return val, err
	}

	bulk := make([]byte, n)

	resp.reader.Read(bulk)

	val.bulkStrVal = strings.ToLower(string(bulk))

	resp.readLine()

	return val, nil
}

// Reads integer
func (resp *Resp) readInteger() (int, int, error) {
	line, n, err := resp.readLine()

	if err != nil {
		return 0, 0, err
	}

	val, err := strconv.ParseInt(string(line), 10, 64)

	if err != nil {
		return 0, 0, err
	}

	return int(val), n, nil
}

// Reads line
func (resp *Resp) readLine() ([]byte, int, error) {
	line := make([]byte, 0)
	n := 0

	for {
		b, err := resp.reader.ReadByte()

		if err != nil {
			return nil, 0, err
		}

		n += 1

		line = append(line, b)

		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}

	return line[:len(line)-2], n, nil
}
