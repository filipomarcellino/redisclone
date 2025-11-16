package main

import "sync"

type KV struct {
	store map[string]string
	lock  sync.RWMutex
}

func NewKV() *KV {
	return &KV{store: make(map[string]string)}
}

func (kv *KV) get(key string) Value {
	defer kv.lock.RUnlock()
	kv.lock.RLock()
	val, ok := kv.store[key]
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: val}
}

func (kv *KV) mget(keys []string) Value {
	defer kv.lock.RUnlock()
	kv.lock.RLock()
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
	defer kv.lock.Unlock()
	kv.lock.Lock()
	for _, pair := range pairs {
		kv.store[pair.key] = pair.value
	}
}

func (kv *KV) set(key string, val string) {
	defer kv.lock.Unlock()
	kv.lock.Lock()
	kv.store[key] = val
}
