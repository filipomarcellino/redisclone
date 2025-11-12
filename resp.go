package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readResp(bufferStr string) string {
	reader := bufio.NewReader(strings.NewReader(bufferStr))

	// read first byte to get type of data
	dataType, _ := reader.ReadByte()
	switch string(dataType) {
	case "+":
		val, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		strings.TrimSuffix(val, "\r")
		return val
	case "$":
		// read second byte to get the length of bulk string
		size, _ := reader.ReadByte()
		strSize, _ := strconv.ParseInt(string(size), 10, 64)
		buffer := make([]byte, strSize)
		reader.Read(buffer)
		return string(buffer)
	case "*":
		fmt.Println("Array")
	// case ":":
	// 	fmt.Println("Integer")
	// case "-":
	// 	fmt.Println("Error")
	default:
		return ""

	}
	return ""
}
