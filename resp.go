package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ParsedType struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []ParsedType
}

type RespParser struct {
	reader *bufio.Reader
}

func newRespParser(r io.Reader) *RespParser {
	return &RespParser{reader: bufio.NewReader(r)}
}

// helper function to read the next integer
func (p *RespParser) readInt() (int, error) {
	sizeStr, _ := p.reader.ReadString('\n')
	trimmed := strings.TrimSuffix(sizeStr, "\r\n")
	size, err := strconv.Atoi(trimmed)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return size, nil
}

func (p *RespParser) readResp() (ParsedType, error) {
	// read first byte to get type of data
	dataType, _ := p.reader.ReadByte()
	switch string(dataType) {
	case "+":
		val, err := p.reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return ParsedType{}, err
		}
		trimmed := strings.TrimSuffix(val, "\r\n")
		return ParsedType{
			typ: "string",
			str: trimmed,
		}, nil
	case "$":
		// read the size of string
		size, err := p.readInt()
		if err != nil {
			return ParsedType{}, err
		}
		// add 2 bytes to consume \r\n at the end of the string
		buffer := make([]byte, size+2)
		p.reader.Read(buffer)
		return ParsedType{
			typ: "bulk",
			str: string(buffer[:size]),
		}, nil
	case "*":
		// read array size
		var parsed ParsedType = ParsedType{}
		parsed.typ = "array"
		size, err := p.readInt()
		if err != nil {
			return ParsedType{}, err
		}
		parsed.array = make([]ParsedType, 0, size)
		for i := 0; i < size; i++ {
			temp, err := p.readResp()
			if err != nil {
				return ParsedType{}, err
			}
			parsed.array = append(parsed.array, temp)
		}
		return parsed, nil

		// case ":":
		// 	fmt.Println("Integer")
		// case "-":
		// 	fmt.Println("Error")
	}
	fmt.Println("Unknown type")
	return ParsedType{}, nil
}
