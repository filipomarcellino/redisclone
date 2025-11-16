package main

import "sync"

type KV struct {
	store map[string]string
	lock  sync.RWMutex
}

func NewKV() *KV {
	return &KV{store: make(map[string]string)}
}

func (kv *KV) get(key string) (string, bool) {
	defer kv.lock.RUnlock()
	kv.lock.RLock()
	val, ok := kv.store[key]
	if !ok {
		return "0", false
	}
	return val, true
}

func (kv *KV) set(key string, val string) {
	defer kv.lock.RUnlock()
	kv.lock.RLock()
	kv.store[key] = val
}
