package main

import "sync"

type KV struct {
	store map[string]any
	lock  sync.RWMutex
}

func NewKV() *KV{
	return &KV{}
}

func (kv *KV) get(key string) (any, bool) {
	defer kv.lock.RUnlock()
	kv.lock.RLock()
	val, ok := kv.store[key]
	if !ok {
		return nil, false
	}
	return val, true
}

func (kv *KV) set(key string, val any) {
	defer kv.lock.RUnlock()
	kv.lock.RLock()
	kv.store[key] = val
}
