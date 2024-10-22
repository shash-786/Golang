package main

import "fmt"

func main() {
	fmt.Println(0)
	// c := make(chan bool)
	go count()
	// fmt.Println(<-c)
	fmt.Println(10)
}

func count() {
	for i := 1; i <= 9; i++ {
		fmt.Println(i)
	}
	// c <- true
}
