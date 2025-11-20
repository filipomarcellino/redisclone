package main

import (
	"strconv"
	"strings"
)

type Executor struct {
	db        []*KV
	currIndex int
	// maybe add pointer to AOF struct here
}

type KeyValuePair struct {
	key   string
	value string
}

func NewExecutor(kvDatabase []*KV) *Executor {
	return &Executor{db: kvDatabase, currIndex: 0}
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
	case "DEL":
		return e.handleDelCommand(input.array[1:])
	case "KEYS":
		return e.handleKeysCommand(input.array[1:])
	case "RENAME":
		return e.handleRenameCommand(input.array[1:])
	case "MSET":
		return e.handleMsetCommand(input.array[1:])
	case "MGET":
		return e.handleMgetCommand(input.array[1:])
	case "SELECT":
		return e.handleSelectCommand(input.array[1:])
	case "FLUSHDB":
		return e.handleFlushDbCommand(input.array[1:])
	case "FLUSHALL":
		return e.handleFlushAllCommand(input.array[1:])
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
	e.db[e.currIndex].set(key, val)
	return Value{typ: "string", str: "OK"}
}

func (e *Executor) handleSetnxCommand(array []Value) Value {
	if len(array) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'setnx' command"}
	}
	key := array[0].bulk
	val := array[1].bulk
	return e.db[e.currIndex].setnx(key, val)
}

func (e *Executor) handleGetCommand(array []Value) Value {
	if len(array) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}
	key := array[0].bulk
	return e.db[e.currIndex].get(key)
}

func (e *Executor) handleDelCommand(array []Value) Value {
	if len(array) < 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'del' command"}
	}
	var keys []string = make([]string, 0, len(array))
	for _, value := range array {
		keys = append(keys, value.bulk)
	}
	return e.db[e.currIndex].del(keys)
}

func (e *Executor) handleKeysCommand(array []Value) Value {
	if len(array) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'keys' command"}
	}
	pattern := array[0].bulk
	return e.db[e.currIndex].keys(pattern)
}

func (e *Executor) handleRenameCommand(array []Value) Value {
	if len(array) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'rename' command"}
	}
	oldKey := array[0].bulk
	newKey := array[1].bulk
	return e.db[e.currIndex].rename(oldKey, newKey)
}

func (e *Executor) handleMsetCommand(array []Value) Value {
	if len(array) < 2 || len(array)%2 == 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'mset' command"}
	}
	var pairs []KeyValuePair = make([]KeyValuePair, 0, len(array)/2)
	for i := 0; i < len(array); i += 2 {
		pairs = append(pairs, KeyValuePair{key: array[i].bulk, value: array[i+1].bulk})
	}
	e.db[e.currIndex].mset(pairs)
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
	return e.db[e.currIndex].mget(keys)
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
	val := e.db[e.currIndex].get(key)

	// new key
	if val.typ == "null" {
		e.db[e.currIndex].set(key, "1")
		return Value{typ: "integer", num: 1}
	}

	// check if value is an integer
	valInt, err := strconv.Atoi(val.bulk)
	if err != nil {
		return Value{typ: "error", str: "ERR value is not an integer or out of range"}
	}

	e.db[e.currIndex].set(key, strconv.Itoa(valInt+1))
	return Value{typ: "integer", num: valInt + 1}
}

func (e *Executor) handleDecrCommand(array []Value) Value {
	if len(array) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'DECR' command"}
	}
	key := array[0].bulk
	val := e.db[e.currIndex].get(key)

	// new key
	if val.typ == "null" {
		e.db[e.currIndex].set(key, "-1")
		return Value{typ: "integer", num: -1}
	}

	// check if value is an integer
	valInt, err := strconv.Atoi(val.bulk)
	if err != nil {
		return Value{typ: "error", str: "ERR value is not an integer or out of range"}
	}

	e.db[e.currIndex].set(key, strconv.Itoa(valInt-1))
	return Value{typ: "integer", num: valInt - 1}
}

func (e *Executor) handleSelectCommand(array []Value) Value {
	if len(array) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'SELECT' command"}
	}
	ind := array[0].bulk
	indInt, err := strconv.Atoi(ind)
	if err != nil {
		return Value{typ: "error", str: "ERR value is not an integer or out of range"}
	}
	if indInt < 0 || indInt > 15 {
		return Value{typ: "error", str: "ERR DB index is out of range"}
	}
	e.currIndex = indInt
	return Value{typ: "string", str: "OK"}
}

func (e *Executor) handleFlushDbCommand(array []Value) Value {
	if len(array) != 0 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'FLUSHDB' command"}
	}
	clear(e.db[e.currIndex].store)
	return Value{typ: "string", str: "OK"}
}

func (e *Executor) handleFlushAllCommand(array []Value) Value {
	if len(array) != 0 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'FLUSHALL' command"}
	}
	for _, kvInstance := range e.db {
		clear(kvInstance.store)
	}
	return Value{typ: "string", str: "OK"}
}
