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
		parser := newRespParser(conn)
		_, err = parser.readResp()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			return
		}
		// parse resp format to string format
		// operation := readResp(string(buffer))
		// res := execute(operation)
		// resRESP := toRESP(res)

		//
		conn.Write([]byte("+OK\r\n"))
	}
}
