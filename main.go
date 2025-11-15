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

	// create unique key value instance
	kv := NewKV()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnetion(conn, kv)

	}
}

func handleConnetion(conn net.Conn, kv *KV) {
	defer conn.Close()
	parser := newRespParser(conn)
	for {
		// kilobyte-size buffer to read messages from client
		val, err := parser.readResp()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			break
		}
		fmt.Printf("request: %+v\n", val)
		executor := NewExecutor(kv)
		responseVal := executor.handleCommand(val)
		fmt.Printf("response: %+v\n", responseVal)
		respBytes := responseVal.Marshal()
		_, err = conn.Write(respBytes)
		if err != nil {
			fmt.Println("error writing to client: ", err.Error())
			break
		}
	}
}
