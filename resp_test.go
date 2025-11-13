package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestSimpleString(t *testing.T) {
	inputSimpleString := "+OK\r\n"
	parser := newRespParser(strings.NewReader(inputSimpleString))
	fmt.Println(parser.readResp())
}

func TestBulkString(t *testing.T) {
	inputBulkString := "$5\r\nAhmed\r\n"
	parser := newRespParser(strings.NewReader(inputBulkString))
	fmt.Println(parser.readResp())
}
func TestArray(t *testing.T) {
	inputArray := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	parser := newRespParser(strings.NewReader(inputArray))
	fmt.Println(parser.readResp())
}
