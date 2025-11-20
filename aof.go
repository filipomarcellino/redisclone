package main

import (
	"os"
	"sync"
)

type AOF struct {
	file *os.File
	lock sync.RWMutex 
}

func newAOF(filename string) {
	// create new file 
}

func (aof *AOF) append(cmd []byte) {
	aof.lock.Lock()
	defer aof.lock.Unlock()
	// append logic here
}