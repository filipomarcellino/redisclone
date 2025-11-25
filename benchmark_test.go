package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	// Initialize KV with 16 shards (for sharded impl) or just NewKV() (for single lock)

	kv := NewKV(128)

	b.ResetTimer()
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := fmt.Sprintf("key-%d", rand.Intn(100000))
			kv.set(key, "value")
		}
	})
}
