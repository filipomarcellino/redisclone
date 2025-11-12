package main

import (
	"bufio"
	"fmt"
	"strings"
)

func readRESP(bufferStr string) {
	// read first byte to know what comes next
	reader := bufio.NewReader(strings.NewReader(bufferStr))

	b, _ := reader.ReadByte()
	fmt.Println(string(b))
	switch string(b) {
	case "+":
		readSimpleString()
	case "$":
		fmt.Println("Bulk string")
	case ":":
		fmt.Println("Integer")
	case "*":
		fmt.Println("Array")
	case "-":
		fmt.Println("Error")
	}
}

func readSimpleString() {
	
}
