package main

import "strings"

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
	case "SET":
		return e.handleSetCommand(input.array[1:])
	case "GET":
		return e.handleGetCommand(input.array[1:])
	default:
		return Value{}
	}
}

func (e *Executor) handleSetCommand(array []Value) Value {
	key := array[0].bulk
	val := array[1].bulk
	e.db.set(key, val)
	return Value{typ: "string", str: "OK"}
}

func (e *Executor) handleGetCommand(array []Value) Value {
	key := array[0].bulk
	val, ok := e.db.get(key)
	valString := val.(string)
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: valString}
}
