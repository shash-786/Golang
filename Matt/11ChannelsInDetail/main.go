package main

import (
	"fmt"
	"time"
)

type T struct {
	i byte
	b bool
}

func do(i int, channel chan<- *T) {
	t := &T{i: byte(i)}
	channel <- t

	// Always first recieve Happens then the send finishes
	// In the case of an unbuffered channel write is blocked until there
	// is an appropriate reader present on the  other end so
	// NOTE: In line number 15 the goroutines will stop until  there
	// is a reader present in arr[i] = *<-channel
	// So by the time the bool value is updated in the do function the read
	// Has Already Happened
	//
	// FIX: To fix this we can make  the channel buffered so that the bool
	// value  can be updated and the goroutines dont block because there isn't a reader
	// present

	t.b = true
}

func main() {
	arr := make([]T, 5)
	// channel := make(chan *T)
	channel := make(chan *T, 5)

	for i := range arr {
		go do(i, channel)
	}

	time.Sleep(1 * time.Second)
	for i := range arr {
		arr[i] = *<-channel
	}

	for _, val := range arr {
		fmt.Println(val)
	}
}
