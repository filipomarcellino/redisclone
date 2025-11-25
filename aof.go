package main

import (
	"os"
	"sync"
)

type AOF struct {
	file *os.File
	lock sync.RWMutex
}

func newAOF(filename string) (*AOF, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &AOF{
		file: file,
	}, nil
}

func (aof *AOF) Close() error {
	aof.lock.Lock()
	defer aof.lock.Unlock()
	return aof.file.Close()
}

func (aof *AOF) append(v Value) {
	aof.lock.Lock()
	defer aof.lock.Unlock()
	rawBytes := v.Marshal()
	_, err := aof.file.Write(rawBytes)
	if err != nil {
		panic(err)
	}
}
