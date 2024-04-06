package main

import (
	"fmt"
	"time"
)

func main() {
	stopper := time.After(10 * time.Second)
	ticker := time.NewTicker(1 * time.Second).C

loop:
	for {
		select {
		case <-ticker:
			fmt.Println("Tick")
		case <-stopper:
			break loop
		}
	}
	fmt.Println("End")
}
