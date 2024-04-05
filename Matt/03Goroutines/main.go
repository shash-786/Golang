package main

import "fmt"

func generate(limit int, write chan<- int) {
	for i := 2; i < limit; i++ {
		write <- i
	}
	close(write)
}

func filter_based_on_prime(prime int, source <-chan int, dst chan<- int) {
	for val := range source {
		if val%prime != 0 {
			dst <- val
		}
	}
	close(dst)
}

func sieve(limit int) {
	channel := make(chan int)
	go generate(limit, channel)

	for {
		prime, ok := <-channel
		if !ok {
			break
		}

		sieved_channel := make(chan int)
		go filter_based_on_prime(prime, channel, sieved_channel)
		channel = sieved_channel
		fmt.Printf("%d ", prime)
	}
	fmt.Println()
}

func main() {
	sieve(150)
}
