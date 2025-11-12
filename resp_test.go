package main

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestSimpleString(t *testing.T) {
	inputSimpleString := "+OK\r\n"
	fmt.Println(readResp(bufio.NewReader(strings.NewReader(inputSimpleString))))
}

func TestBulkString(t *testing.T) {
	inputBulkString := "$5\r\nAhmed\r\n"
	fmt.Println(readResp(bufio.NewReader(strings.NewReader(inputBulkString))))
}
func TestArray(t *testing.T) {
	inputArray := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	fmt.Println(readResp(bufio.NewReader(strings.NewReader(inputArray))))
}
