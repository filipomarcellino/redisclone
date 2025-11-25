package main

import (
	"strconv"
	"strings"
)

type Executor struct {
	db  *KV
	aof *AOF
}

type KeyValuePair struct {
	key   string
	value string
}

func NewExecutor(kvDatabase *KV, aofPointer *AOF) *Executor {
	return &Executor{db: kvDatabase, aof: aofPointer}
}

func (e *Executor) handleCommand(input Value) Value {
	// always an array type
	if input.typ != "array" {
		return Value{typ: "error", str: "ERR expected array type"}
	}
	switch strings.ToUpper(input.array[0].bulk) {
	case "PING":
		return e.handlePingCommand(input.array[1:])
	case "INCR":
		res := e.handleIncrCommand(input.array[1:])
		if res.typ != "error" {
			e.persistToAOF(input)
		}
		return res
	case "DECR":
		res := e.handleDecrCommand(input.array[1:])
		if res.typ != "error" {
			e.persistToAOF(input)
		}
		return res
	case "SET":
		res := e.handleSetCommand(input.array[1:])
		if res.typ != "error" {
			e.persistToAOF(input)
		}
		return res
	case "SETNX":
		res := e.handleSetnxCommand(input.array[1:])
		if res.typ != "error" {
			e.persistToAOF(input)
		}
		return res
	case "GET":
		return e.handleGetCommand(input.array[1:])
	case "DEL":
		res := e.handleDelCommand(input.array[1:])
		if res.typ != "error" {
			e.persistToAOF(input)
		}
		return res
	case "KEYS":
		return e.handleKeysCommand(input.array[1:])
	case "RENAME":
		res := e.handleRenameCommand(input.array[1:])
		if res.typ != "error" {
			e.persistToAOF(input)
		}
		return res
	case "MSET":
		res := e.handleMsetCommand(input.array[1:])
		if res.typ != "error" {
			e.persistToAOF(input)
		}
		return res
	case "MGET":
		return e.handleMgetCommand(input.array[1:])
	case "FLUSHDB":
		res := e.handleFlushDbCommand(input.array[1:])
		if res.typ != "error" {
			e.persistToAOF(input)
		}
		return res
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

func (e *Executor) persistToAOF(v Value) {
	if e.aof != nil {
		e.aof.append(v)
	}
}

func (e *Executor) handleSetCommand(array []Value) Value {
	if len(array) < 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}
	key := array[0].bulk
	val := array[1].bulk
	e.db.set(key, val)
	res := Value{typ: "string", str: "OK"}
	return res
}

func (e *Executor) handleSetnxCommand(array []Value) Value {
	if len(array) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'setnx' command"}
	}
	key := array[0].bulk
	val := array[1].bulk
	res := e.db.setnx(key, val)
	return res
}

func (e *Executor) handleGetCommand(array []Value) Value {
	if len(array) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}
	key := array[0].bulk
	return e.db.get(key)
}

func (e *Executor) handleDelCommand(array []Value) Value {
	if len(array) < 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'del' command"}
	}
	var keys []string = make([]string, 0, len(array))
	for _, value := range array {
		keys = append(keys, value.bulk)
	}
	res := e.db.del(keys)
	return res
}

func (e *Executor) handleKeysCommand(array []Value) Value {
	if len(array) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'keys' command"}
	}
	pattern := array[0].bulk
	return e.db.keys(pattern)
}

func (e *Executor) handleRenameCommand(array []Value) Value {
	if len(array) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'rename' command"}
	}
	oldKey := array[0].bulk
	newKey := array[1].bulk
	res := e.db.rename(oldKey, newKey)
	return res
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
	res := Value{typ: "string", str: "OK"}
	return res
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
	res := Value{typ: "integer", num: valInt + 1}
	return res
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
		return Value{typ: "integer", num: -1}
	}

	// check if value is an integer
	valInt, err := strconv.Atoi(val.bulk)
	if err != nil {
		return Value{typ: "error", str: "ERR value is not an integer or out of range"}
	}

	e.db.set(key, strconv.Itoa(valInt-1))
	res := Value{typ: "integer", num: valInt - 1}
	return res
}

func (e *Executor) handleFlushDbCommand(array []Value) Value {
	if len(array) != 0 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'FLUSHDB' command"}
	}
	e.db.Flush()
	res := Value{typ: "string", str: "OK"}
	return res
}
