package main

import (
	"path"
	"sync"
)

type KV struct {
	store map[string]string
	lock  sync.RWMutex
}

func NewKV() *KV {
	return &KV{store: make(map[string]string)}
}

func (kv *KV) get(key string) Value {
	kv.lock.RLock()
	defer kv.lock.RUnlock()
	val, ok := kv.store[key]
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: val}
}

func (kv *KV) mget(keys []string) Value {
	kv.lock.RLock()
	defer kv.lock.RUnlock()
	res := Value{}
	res.typ = "array"
	for _, key := range keys {
		val, ok := kv.store[key]
		if !ok {
			res.array = append(res.array, Value{typ: "null"})
			continue
		}
		res.array = append(res.array, Value{typ: "bulk", bulk: val})
	}
	return res
}

// atomic multiple set operation
func (kv *KV) mset(pairs []KeyValuePair) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	for _, pair := range pairs {
		kv.store[pair.key] = pair.value
	}
}

func (kv *KV) set(key string, val string) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	kv.store[key] = val
}

func (kv *KV) setnx(key string, val string) Value {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	_, ok := kv.store[key]
	// early return when data already exists
	if ok {
		return Value{typ: "integer", num: 0}
	}
	kv.store[key] = val
	return Value{typ: "integer", num: 1}
}

func (kv *KV) del(keys []string) Value {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	count := 0
	for _, key := range keys {
		_, ok := kv.store[key]
		if ok {
			delete(kv.store, key)
			count++
		}
	}
	return Value{typ: "integer", num: count}
}

func (kv *KV) keys(pattern string) Value {
	kv.lock.RLock()
	defer kv.lock.RUnlock()
	res := Value{}
	res.typ = "array"
	res.array = []Value{}
	for key := range kv.store {
		matched, err := path.Match(pattern, key)
		if err != nil {
			return Value{typ: "error", str: "ERR invalid pattern"}
		}
		if matched {
			res.array = append(res.array, Value{typ: "bulk", bulk: key})
		}
	}
	return res
}

func (kv *KV) rename(oldKey string, newKey string) Value {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	val, ok := kv.store[oldKey]
	if !ok {
		return Value{typ: "error", str: "ERR no such key"}
	}
	kv.store[newKey] = val
	delete(kv.store, oldKey)
	return Value{typ: "string", str: "OK"}
}