package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	STRING byte = '+'
	BULK   byte = '$'
	ARRAY  byte = '*'
	ERROR  byte = '-'
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

func (pt ParsedType) Marshal() []byte {
	switch pt.typ {
	case "array":
		return pt.MarshalArray()
	case "bulk":
		return pt.MarshalBulk()
	case "string":
		return pt.MarshalString()
	case "null":
		return pt.marshallNull()
	case "error":
		return pt.marshallError()
	default:
		return []byte{}
	}
}

func (pt ParsedType) MarshalString() []byte {
	buffer := []byte{}
	buffer = append(buffer, STRING)
	buffer = append(buffer, pt.str...)
	buffer = append(buffer, '\r', '\n')
	return buffer
}

func (pt ParsedType) MarshalBulk() []byte {
	buffer := []byte{}
	buffer = append(buffer, BULK)
	buffer = append(buffer, strconv.Itoa(len(pt.bulk))...)
	buffer = append(buffer, '\r', '\n')
	buffer = append(buffer, pt.bulk...)
	buffer = append(buffer, '\r', '\n')
	return buffer
}

func (pt ParsedType) MarshalArray() []byte {
	buffer := []byte{}
	buffer = append(buffer, ARRAY)
	buffer = append(buffer, strconv.Itoa(len(pt.array))...)
	buffer = append(buffer, '\r', '\n')

	// loop through array
	for _, val := range pt.array {
		buffer = append(buffer, val.Marshal()...)
	}
	return buffer
}
func (pt ParsedType) marshallError() []byte {
	var buffer []byte
	buffer = append(buffer, ERROR)
	buffer = append(buffer, pt.str...)
	buffer = append(buffer, '\r', '\n')

	return buffer
}

func (pt ParsedType) marshallNull() []byte {
	return []byte("$-1\r\n")
}
func newRespParser(rd io.Reader) *RespParser {
	return &RespParser{reader: bufio.NewReader(rd)}
}

// helper function to read the next integer
func (r *RespParser) readInt() (int, error) {
	sizeStr, _ := r.reader.ReadString('\n')
	trimmed := strings.TrimSuffix(sizeStr, "\r\n")
	size, err := strconv.Atoi(trimmed)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return size, nil
}

func (r *RespParser) readResp() (ParsedType, error) {
	// read first byte to get type of data
	dataType, err := r.reader.ReadByte()
	if err != nil {
		return ParsedType{}, err
	}
	switch dataType {
	case STRING:
		val, err := r.reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return ParsedType{}, err
		}
		trimmed := strings.TrimSuffix(val, "\r\n")
		return ParsedType{
			typ: "string",
			str: trimmed,
		}, nil
	case BULK:
		// read the size of string
		size, err := r.readInt()
		if err != nil {
			return ParsedType{}, err
		}
		// add 2 bytes to consume \r\n at the end of the string
		buffer := make([]byte, size+2)
		r.reader.Read(buffer)
		return ParsedType{
			typ:  "bulk",
			bulk: string(buffer[:size]),
		}, nil
	case ARRAY:
		// read array size
		var parsed ParsedType = ParsedType{}
		parsed.typ = "array"
		size, err := r.readInt()
		if err != nil {
			return ParsedType{}, err
		}
		parsed.array = make([]ParsedType, 0, size)
		for i := 0; i < size; i++ {
			temp, err := r.readResp()
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
	default:
		return ParsedType{}, fmt.Errorf("unknown RESP type: %c (byte: %d)", dataType, dataType)
	}
}
