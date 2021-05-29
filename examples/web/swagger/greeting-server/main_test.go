package main

import (
	"sync"
	"testing"
)

var mu sync.Mutex
func TestRunMain(t *testing.T) {
	mu.Lock()
	go main()
	mu.Unlock()
}
