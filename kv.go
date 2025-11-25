package main

import (
	"hash/fnv"
	"path"
	"sync"
)

type KV struct {
	shards     []*Shard
	shardCount int
}

type Shard struct {
	store map[string]string
	lock  sync.RWMutex
	id    int
}

func NewKV(shardCount int) *KV {
	shards := make([]*Shard, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = &Shard{store: make(map[string]string), id: i}
	}
	return &KV{shards: shards, shardCount: shardCount}
}

func (kv *KV) getShard(key string) *Shard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return kv.shards[int(h.Sum32())%kv.shardCount]
}

func (kv *KV) get(key string) Value {
	shard := kv.getShard(key)
	shard.lock.RLock()
	defer shard.lock.RUnlock()
	val, ok := shard.store[key]
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: val}
}

func (kv *KV) mget(keys []string) Value {
	res := Value{}
	res.typ = "array"
	for _, key := range keys {
		shard := kv.getShard(key)
		shard.lock.RLock()
		val, ok := shard.store[key]
		shard.lock.RUnlock()
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
	// group by shard
	shardUpdates := make(map[*Shard][]KeyValuePair)
	for _, pair := range pairs {
		shard := kv.getShard(pair.key)
		shardUpdates[shard] = append(shardUpdates[shard], pair)
	}

	// apply updates per shard
	for shard, updates := range shardUpdates {
		shard.lock.Lock()
		for _, pair := range updates {
			shard.store[pair.key] = pair.value
		}
		shard.lock.Unlock()
	}
}

func (kv *KV) set(key string, val string) {
	shard := kv.getShard(key)
	shard.lock.Lock()
	defer shard.lock.Unlock()
	shard.store[key] = val
}

func (kv *KV) setnx(key string, val string) Value {
	shard := kv.getShard(key)
	shard.lock.Lock()
	defer shard.lock.Unlock()
	_, ok := shard.store[key]
	// early return when data already exists
	if ok {
		return Value{typ: "integer", num: 0}
	}
	shard.store[key] = val
	return Value{typ: "integer", num: 1}
}

func (kv *KV) del(keys []string) Value {
	// group by shard
	shardKeys := make(map[*Shard][]string)
	for _, key := range keys {
		shard := kv.getShard(key)
		shardKeys[shard] = append(shardKeys[shard], key)
	}

	count := 0
	// apply deletes per shard
	for shard, keys := range shardKeys {
		shard.lock.Lock()
		for _, key := range keys {
			_, ok := shard.store[key]
			if ok {
				delete(shard.store, key)
				count++
			}
		}
		shard.lock.Unlock()
	}
	return Value{typ: "integer", num: count}
}

func (kv *KV) keys(pattern string) Value {
	res := Value{}
	res.typ = "array"
	res.array = []Value{}
	for _, shard := range kv.shards {
		shard.lock.RLock()
		for key := range shard.store {
			matched, err := path.Match(pattern, key)
			if err != nil {
				shard.lock.RUnlock()
				return Value{typ: "error", str: "ERR invalid pattern"}
			}
			if matched {
				res.array = append(res.array, Value{typ: "bulk", bulk: key})
			}
		}
		shard.lock.RUnlock()
	}
	return res
}

func (kv *KV) rename(oldKey string, newKey string) Value {
	oldShard := kv.getShard(oldKey)
	newShard := kv.getShard(newKey)

	// Same shard optimization
	if oldShard == newShard {
		oldShard.lock.Lock()
		defer oldShard.lock.Unlock()
		val, ok := oldShard.store[oldKey]
		if !ok {
			return Value{typ: "error", str: "ERR no such key"}
		}
		oldShard.store[newKey] = val
		delete(oldShard.store, oldKey)
		return Value{typ: "string", str: "OK"}
	}

	// Different shards: lock in order to avoid deadlock
	firstLock := oldShard
	secondLock := newShard
	if oldShard.id > newShard.id {
		firstLock = newShard
		secondLock = oldShard
	}

	firstLock.lock.Lock()
	defer firstLock.lock.Unlock()
	secondLock.lock.Lock()
	defer secondLock.lock.Unlock()

	val, ok := oldShard.store[oldKey]
	if !ok {
		return Value{typ: "error", str: "ERR no such key"}
	}

	newShard.store[newKey] = val
	delete(oldShard.store, oldKey)

	return Value{typ: "string", str: "OK"}
}

func (kv *KV) Flush() {
	for _, shard := range kv.shards {
		shard.lock.Lock()
		clear(shard.store)
		shard.lock.Unlock()
	}
}
