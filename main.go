package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	fmt.Println("Listening on port :6379")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()
	for {
		// kilobyte-size buffer to read messages from client
		buffer := make([]byte, 1024)

		_, err = conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			return
		}
		readRESP(string(buffer))
		conn.Write([]byte("+OK\r\n"))
	}
}
