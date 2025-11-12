package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readResp(reader *bufio.Reader) []string {
	// read first byte to get type of data
	dataType, _ := reader.ReadByte()
	switch string(dataType) {
	case "+":
		val, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		trimmed := strings.TrimSuffix(val, "\r\n")
		return []string{trimmed}
	case "$":
		// read the size of string
		sizeStr, _ := reader.ReadString('\n')
		trimmed := strings.TrimSuffix(sizeStr, "\r\n")
		size, err := strconv.Atoi(trimmed)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// add 2 bytes to consume \r\n at the end of the string
		buffer := make([]byte, size + 2)
		reader.Read(buffer)
		return []string {string(buffer[: size])}
	case "*":
		sizeStr, _ := reader.ReadString('\n')
		trimmed := strings.TrimSuffix(sizeStr, "\r\n")
		size, err := strconv.Atoi(trimmed)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		var slice []string = make([]string, 0, size)
		
		for i := 0; i < size; i++ {
			slice = append(slice, readResp(reader)...)
		}
		return slice

	// case ":":
	// 	fmt.Println("Integer")
	// case "-":
	// 	fmt.Println("Error")
	default:
		return nil

	}
	return nil
}
