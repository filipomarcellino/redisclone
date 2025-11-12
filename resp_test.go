package main

import (
	"testing"
)

func TestReadRESP(t *testing.T) {
	input := "$5\r\nAhmed\r\n"
	read(input)

}
