package main

import (
	"fmt"
	"sync"
	"time"
)

type shared_variant struct {
	data int
	mu   sync.Mutex
}

func main() {
	s := &shared_variant{
		data: 0,
	}

	for i := 0; i < 5; i++ {
		go func(s *shared_variant) {
			s.mu.Lock()
			fmt.Printf("before: %d\n", s.data)
			s.data++
			fmt.Printf("after: %d\n", s.data)
			s.mu.Unlock()
		}(s)
	}

	time.Sleep(2 * time.Second)
}
