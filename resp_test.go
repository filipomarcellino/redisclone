package main

import (
	"bytes"
	"strings"
	"testing"
)

// Tests for RespParser.readResp() method

func TestSimpleString(t *testing.T) {
	inputSimpleString := "+OK\r\n"
	parser := newRespParser(strings.NewReader(inputSimpleString))
	result, err := parser.readResp()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.typ != "string" {
		t.Errorf("Expected type 'string', got '%s'", result.typ)
	}
	if result.str != "OK" {
		t.Errorf("Expected 'OK', got '%s'", result.str)
	}
}

func TestSimpleStringEmpty(t *testing.T) {
	inputSimpleString := "+\r\n"
	parser := newRespParser(strings.NewReader(inputSimpleString))
	result, err := parser.readResp()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.typ != "string" {
		t.Errorf("Expected type 'string', got '%s'", result.typ)
	}
	if result.str != "" {
		t.Errorf("Expected empty string, got '%s'", result.str)
	}
}

func TestSimpleStringWithSpaces(t *testing.T) {
	inputSimpleString := "+Hello World\r\n"
	parser := newRespParser(strings.NewReader(inputSimpleString))
	result, err := parser.readResp()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.str != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result.str)
	}
}

func TestBulkString(t *testing.T) {
	inputBulkString := "$5\r\nAhmed\r\n"
	parser := newRespParser(strings.NewReader(inputBulkString))
	result, err := parser.readResp()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.typ != "bulk" {
		t.Errorf("Expected type 'bulk', got '%s'", result.typ)
	}
	if result.bulk != "Ahmed" {
		t.Errorf("Expected 'Ahmed', got '%s'", result.bulk)
	}
}

func TestBulkStringEmpty(t *testing.T) {
	inputBulkString := "$0\r\n\r\n"
	parser := newRespParser(strings.NewReader(inputBulkString))
	result, err := parser.readResp()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.typ != "bulk" {
		t.Errorf("Expected type 'bulk', got '%s'", result.typ)
	}
	if result.bulk != "" {
		t.Errorf("Expected empty string, got '%s'", result.bulk)
	}
}

func TestBulkStringWithSpecialChars(t *testing.T) {
	inputBulkString := "$13\r\nHello\r\nWorld!\r\n"
	parser := newRespParser(strings.NewReader(inputBulkString))
	result, err := parser.readResp()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.bulk != "Hello\r\nWorld!" {
		t.Errorf("Expected 'Hello\\r\\nWorld!', got '%s'", result.bulk)
	}
}

func TestArray(t *testing.T) {
	inputArray := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	parser := newRespParser(strings.NewReader(inputArray))
	result, err := parser.readResp()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.typ != "array" {
		t.Errorf("Expected type 'array', got '%s'", result.typ)
	}
	if len(result.array) != 2 {
		t.Errorf("Expected array length 2, got %d", len(result.array))
	}
	if result.array[0].bulk != "hello" {
		t.Errorf("Expected first element 'hello', got '%s'", result.array[0].bulk)
	}
	if result.array[1].bulk != "world" {
		t.Errorf("Expected second element 'world', got '%s'", result.array[1].bulk)
	}
}

func TestArrayEmpty(t *testing.T) {
	inputArray := "*0\r\n"
	parser := newRespParser(strings.NewReader(inputArray))
	result, err := parser.readResp()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.typ != "array" {
		t.Errorf("Expected type 'array', got '%s'", result.typ)
	}
	if len(result.array) != 0 {
		t.Errorf("Expected empty array, got length %d", len(result.array))
	}
}

func TestArrayMixed(t *testing.T) {
	inputArray := "*3\r\n+simple\r\n$4\r\nbulk\r\n+another\r\n"
	parser := newRespParser(strings.NewReader(inputArray))
	result, err := parser.readResp()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result.array) != 3 {
		t.Errorf("Expected array length 3, got %d", len(result.array))
	}
	if result.array[0].typ != "string" || result.array[0].str != "simple" {
		t.Errorf("Expected first element to be simple string 'simple'")
	}
	if result.array[1].typ != "bulk" || result.array[1].bulk != "bulk" {
		t.Errorf("Expected second element to be bulk string 'bulk'")
	}
	if result.array[2].typ != "string" || result.array[2].str != "another" {
		t.Errorf("Expected third element to be simple string 'another'")
	}
}

func TestArrayNested(t *testing.T) {
	inputArray := "*2\r\n*2\r\n+inner1\r\n+inner2\r\n+outer\r\n"
	parser := newRespParser(strings.NewReader(inputArray))
	result, err := parser.readResp()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result.array) != 2 {
		t.Errorf("Expected outer array length 2, got %d", len(result.array))
	}
	if result.array[0].typ != "array" {
		t.Errorf("Expected first element to be array")
	}
	if len(result.array[0].array) != 2 {
		t.Errorf("Expected nested array length 2, got %d", len(result.array[0].array))
	}
}

func TestReadIntValid(t *testing.T) {
	parser := newRespParser(strings.NewReader("42\r\n"))
	result, err := parser.readInt()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestReadIntZero(t *testing.T) {
	parser := newRespParser(strings.NewReader("0\r\n"))
	result, err := parser.readInt()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestReadIntNegative(t *testing.T) {
	parser := newRespParser(strings.NewReader("-1\r\n"))
	result, err := parser.readInt()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != -1 {
		t.Errorf("Expected -1, got %d", result)
	}
}

func TestReadIntInvalid(t *testing.T) {
	parser := newRespParser(strings.NewReader("abc\r\n"))
	_, err := parser.readInt()

	if err == nil {
		t.Error("Expected error for invalid integer, got nil")
	}
}

// Tests for ParsedType Marshal methods

func TestMarshalString(t *testing.T) {
	pt := Value{
		typ: "string",
		str: "OK",
	}
	result := pt.MarshalString()
	expected := []byte("+OK\r\n")

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMarshalStringEmpty(t *testing.T) {
	pt := Value{
		typ: "string",
		str: "",
	}
	result := pt.MarshalString()
	expected := []byte("+\r\n")

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMarshalStringWithSpaces(t *testing.T) {
	pt := Value{
		typ: "string",
		str: "Hello World",
	}
	result := pt.MarshalString()
	expected := []byte("+Hello World\r\n")

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMarshalBulk(t *testing.T) {
	pt := Value{
		typ:  "bulk",
		bulk: "hello",
	}
	result := pt.MarshalBulk()
	expected := []byte("$5\r\nhello\r\n")

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMarshalBulkEmpty(t *testing.T) {
	pt := Value{
		typ:  "bulk",
		bulk: "",
	}
	result := pt.MarshalBulk()
	expected := []byte("$0\r\n\r\n")

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMarshallError(t *testing.T) {
	pt := Value{
		typ: "error",
		str: "ERR unknown command",
	}
	result := pt.marshallError()
	expected := []byte("-ERR unknown command\r\n")

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMarshallErrorEmpty(t *testing.T) {
	pt := Value{
		typ: "error",
		str: "",
	}
	result := pt.marshallError()
	expected := []byte("-\r\n")

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMarshallNull(t *testing.T) {
	pt := Value{
		typ: "null",
	}
	result := pt.marshallNull()
	expected := []byte("$-1\r\n")

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMarshalArray(t *testing.T) {
	pt := Value{
		typ: "array",
		array: []Value{
			{typ: "string", str: "hello"},
			{typ: "string", str: "world"},
		},
	}
	result := pt.MarshalArray()
	expected := []byte("*2\r\n+hello\r\n+world\r\n")

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMarshalArrayEmpty(t *testing.T) {
	pt := Value{
		typ:   "array",
		array: []Value{},
	}
	result := pt.MarshalArray()
	expected := []byte("*0\r\n")

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMarshalDefault(t *testing.T) {
	pt := Value{
		typ: "unknown",
	}
	result := pt.Marshal()
	expected := []byte{}

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected empty byte slice, got %v", result)
	}
}

func TestMarshalNull(t *testing.T) {
	pt := Value{
		typ: "null",
	}
	result := pt.Marshal()
	expected := []byte("$-1\r\n")

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// Edge case tests

func TestUnknownDataType(t *testing.T) {
	inputUnknown := ":123\r\n"
	parser := newRespParser(strings.NewReader(inputUnknown))
	_, err := parser.readResp()

	// Should return an error for unknown RESP type
	if err == nil {
		t.Error("Expected error for unknown data type, got nil")
	}
}

func TestNewRespParser(t *testing.T) {
	reader := strings.NewReader("test")
	parser := newRespParser(reader)

	if parser == nil {
		t.Error("Expected non-nil parser")
	}
	if parser.reader == nil {
		t.Error("Expected non-nil reader in parser")
	}
}

// Integration test: parse and then marshal
func TestParseAndMarshalSimpleString(t *testing.T) {
	input := "+Hello\r\n"
	parser := newRespParser(strings.NewReader(input))
	parsed, err := parser.readResp()

	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	marshalled := parsed.MarshalString()
	expected := []byte("+Hello\r\n")

	if !bytes.Equal(marshalled, expected) {
		t.Errorf("Expected %v, got %v", expected, marshalled)
	}
}
