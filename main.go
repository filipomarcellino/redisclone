package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Listening on port :6379")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// create 16 kv instances
	var kvDatabase []*KV = make([]*KV, 16)
	for i := range 16 {
		kvDatabase[i] = NewKV()
	}
	// todo: search for aof file in the enclosing directory
	// scanForAof()
	// todo: load aof file 
	// loadAOF(file, kvDatabase)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(conn, kvDatabase)

	}
}

func loadAOF(file os.File, kvDatabase []*KV) {
	// aof := newAOF(file)
	// aofParser := newRespParser(aof.file)
	executor := NewExecutor(kvDatabase)
	for {
		// kilobyte-size buffer to read messages from client
		val, err := aofParser.readResp()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			break
		}
		executor.handleCommand(val)
	}
}
func handleConnection(conn net.Conn, kvDatabase []*KV) {
	defer conn.Close()
	parser := newRespParser(conn)
	executor := NewExecutor(kvDatabase)
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
		responseVal := executor.handleCommand(val)
		fmt.Printf("response: %+v\n", responseVal)
		if responseVal.typ == "quit" {
			break
		}
		respBytes := responseVal.Marshal()
		_, err = conn.Write(respBytes)
		if err != nil {
			fmt.Println("error writing to client: ", err.Error())
			break
		}
	}
}
