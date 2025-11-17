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
	INT    byte = ':'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type RespParser struct {
	reader *bufio.Reader
}

func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.MarshalArray()
	case "bulk":
		return v.MarshalBulk()
	case "string":
		return v.MarshalString()
	case "integer":
		return v.MarshalInt()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

func (v Value) MarshalString() []byte {
	buffer := []byte{}
	buffer = append(buffer, STRING)
	buffer = append(buffer, v.str...)
	buffer = append(buffer, '\r', '\n')
	return buffer
}

func (v Value) MarshalInt() []byte {
	buffer := []byte{}
	buffer = append(buffer, INT)
	buffer = append(buffer, strconv.Itoa(v.num)...)
	buffer = append(buffer, '\r', '\n')
	return buffer
}

func (v Value) MarshalBulk() []byte {
	buffer := []byte{}
	buffer = append(buffer, BULK)
	buffer = append(buffer, strconv.Itoa(len(v.bulk))...)
	buffer = append(buffer, '\r', '\n')
	buffer = append(buffer, v.bulk...)
	buffer = append(buffer, '\r', '\n')
	return buffer
}

func (v Value) MarshalArray() []byte {
	buffer := []byte{}
	buffer = append(buffer, ARRAY)
	buffer = append(buffer, strconv.Itoa(len(v.array))...)
	buffer = append(buffer, '\r', '\n')

	// loop through array
	for _, val := range v.array {
		buffer = append(buffer, val.Marshal()...)
	}
	return buffer
}
func (v Value) marshalError() []byte {
	var buffer []byte
	buffer = append(buffer, ERROR)
	buffer = append(buffer, v.str...)
	buffer = append(buffer, '\r', '\n')

	return buffer
}

func (v Value) marshalNull() []byte {
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

func (r *RespParser) readResp() (Value, error) {
	// read first byte to get type of data
	dataType, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch dataType {
	case STRING:
		val, err := r.reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return Value{}, err
		}
		trimmed := strings.TrimSuffix(val, "\r\n")
		return Value{
			typ: "string",
			str: trimmed,
		}, nil
	case INT:
		val, err := r.readInt()
		if err != nil {
			return Value{}, err
		}
		return Value{
			typ: "integer",
			num: val,
		}, nil

	case BULK:
		// read the size of string
		size, err := r.readInt()
		if err != nil {
			return Value{}, err
		}
		// add 2 bytes to consume \r\n at the end of the string
		buffer := make([]byte, size+2)
		r.reader.Read(buffer)
		return Value{
			typ:  "bulk",
			bulk: string(buffer[:size]),
		}, nil
	case ARRAY:
		// read array size
		var parsed Value = Value{}
		parsed.typ = "array"
		size, err := r.readInt()
		if err != nil {
			return Value{}, err
		}
		parsed.array = make([]Value, 0, size)
		for i := 0; i < size; i++ {
			temp, err := r.readResp()
			if err != nil {
				return Value{}, err
			}
			parsed.array = append(parsed.array, temp)
		}
		return parsed, nil

	// case ":":
	// 	fmt.Println("Integer")
	// case "-":
	// 	fmt.Println("Error")
	default:
		return Value{}, fmt.Errorf("unknown RESP type: %c (byte: %d)", dataType, dataType)
	}
}
