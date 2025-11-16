package main

import (
	"os"
	"strings"
)

type Executor struct {
	db *KV
}

func NewExecutor(db *KV) *Executor {
	return &Executor{db: db}
}

func (e *Executor) handleCommand(input Value) Value {
	// always an array type
	if input.typ != "array" {
		//todo: return err
	}
	switch strings.ToUpper(input.array[0].bulk) {
	case "PING":
		return e.handlePingCommand(input.array[1:])
	case "QUIT":
		os.Exit(1)
		return Value{}
	case "SET":
		return e.handleSetCommand(input.array[1:])
	case "GET":
		return e.handleGetCommand(input.array[1:])
	case "COMMAND":
		// redis-cli asks for "COMMAND DOCS" or just "COMMAND" on startup for smart auto-completion
		// we'll stub this implementation for now by returning an empty array
		return Value{typ: "array", array: []Value{}}
	default:
		return Value{}
	}
}

func (e *Executor) handleSetCommand(array []Value) Value {
	if len(array) < 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}
	key := array[0].bulk
	val := array[1].bulk
	e.db.set(key, val)
	return Value{typ: "string", str: "OK"}
}

func (e *Executor) handleGetCommand(array []Value) Value {
	if len(array) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}
	key := array[0].bulk
	val, ok := e.db.get(key)
	valString := val.(string)
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: valString}
}

func (e *Executor) handlePingCommand(array []Value) Value {
	switch len(array) {
	case 0:
		return Value{typ: "string", str: "PONG"}
	case 1:
		return Value{typ: "string", str: array[0].bulk}
	default:
		return Value{typ: "error", str: "ERR wrong number of arguments for 'ping' command"}
	}
}
