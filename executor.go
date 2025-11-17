package main

import (
	"strconv"
	"strings"
)

type Executor struct {
	db *KV
}

type KeyValuePair struct {
	key   string
	value string
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
	case "INCR":
		return e.handleIncrCommand(input.array[1:])
	case "DECR":
		return e.handleDecrCommand(input.array[1:])
	case "SET":
		return e.handleSetCommand(input.array[1:])
	case "SETNX":
		return e.handleSetnxCommand(input.array[1:])
	case "GET":
		return e.handleGetCommand(input.array[1:])
	case "MSET":
		return e.handleMsetCommand(input.array[1:])
	case "MGET":
		return e.handleMgetCommand(input.array[1:])
	case "COMMAND":
		// redis-cli asks for "COMMAND DOCS" or just "COMMAND" on startup for smart auto-completion
		// we'll stub this implementation for now by returning an empty array
		return Value{typ: "array", array: []Value{}}
	case "QUIT":
		return Value{typ: "quit"}
	default:
		return Value{typ: "error", str: "ERR command not implemented yet"}
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

func (e *Executor) handleSetnxCommand(array []Value) Value {
	if len(array) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'setnx' command"}
	}
	key := array[0].bulk
	val := array[1].bulk
	return e.db.setnx(key, val)
}

func (e *Executor) handleGetCommand(array []Value) Value {
	if len(array) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}
	key := array[0].bulk
	return e.db.get(key)
}

func (e *Executor) handleMsetCommand(array []Value) Value {
	if len(array) < 2 || len(array)%2 == 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'mset' command"}
	}
	var pairs []KeyValuePair = make([]KeyValuePair, 0, len(array)/2)
	for i := 0; i < len(array); i += 2 {
		pairs = append(pairs, KeyValuePair{key: array[i].bulk, value: array[i+1].bulk})
	}
	e.db.mset(pairs)
	return Value{typ: "string", str: "OK"}
}

func (e *Executor) handleMgetCommand(array []Value) Value {
	if len(array) < 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'mget' command"}
	}
	var keys []string = make([]string, 0, len(array))
	for _, value := range array {
		keys = append(keys, value.bulk)
	}
	return e.db.mget(keys)
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

func (e *Executor) handleIncrCommand(array []Value) Value {
	if len(array) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'INCR' command"}
	}
	key := array[0].bulk
	val := e.db.get(key)

	// new key
	if val.typ == "null" {
		e.db.set(key, "1")
		return Value{typ: "integer", num: 1}
	}

	// check if value is an integer
	valInt, err := strconv.Atoi(val.bulk)
	if err != nil {
		return Value{typ: "error", str: "ERR value is not an integer or out of range"}
	}

	e.db.set(key, strconv.Itoa(valInt+1))
	return Value{typ: "integer", num: valInt + 1}
}

func (e *Executor) handleDecrCommand(array []Value) Value {
	if len(array) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'DECR' command"}
	}
	key := array[0].bulk
	val := e.db.get(key)

	// new key
	if val.typ == "null" {
		e.db.set(key, "-1")
		return Value{typ: "integer", num: 1}
	}

	// check if value is an integer
	valInt, err := strconv.Atoi(val.bulk)
	if err != nil {
		return Value{typ: "error", str: "ERR value is not an integer or out of range"}
	}

	e.db.set(key, strconv.Itoa(valInt-1))
	return Value{typ: "integer", num: valInt - 1}
}
