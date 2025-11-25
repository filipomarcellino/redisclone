package main

import (
	"bufio"
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

	// create a single kv instance
	kvDatabase := NewKV()

	// initialize AOF
	aof, err := newAOF("append-only.aof")
	if err != nil {
		fmt.Println("error initializing AOF:", err)
		return
	}
	defer aof.Close()

	// load AOF file if it exists
	loadAOF(kvDatabase, aof)
	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				fmt.Printf("accept temp error: %v\n", err)
				continue
			}

			fmt.Printf("accept error: %v\n", err)
			return
		}
		go handleConnection(conn, kvDatabase, aof)
	}
}

func loadAOF(kvDatabase *KV, aof *AOF) {
	aofParser := newRespParser(aof.file)
	// pass aof pointer as nil because we don't want to write to aof while reading from it
	executor := NewExecutor(kvDatabase, nil)
	for {
		val, err := aofParser.readResp()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from AOF: ", err.Error())
			break
		}
		executor.handleCommand(val)
	}
}

func handleConnection(conn net.Conn, kvDatabase *KV, aof *AOF) {
	defer conn.Close()
	parser := newRespParser(conn)
	executor := NewExecutor(kvDatabase, aof)
	writer := bufio.NewWriter(conn)
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
		responseVal := executor.handleCommand(val)
		if responseVal.typ == "quit" {
			break
		}
		respBytes := responseVal.Marshal()
		_, err = conn.Write(respBytes)
		if err != nil {
			fmt.Println("error writing to client: ", err.Error())
			break
		}
		writer.Flush()
	}
}
