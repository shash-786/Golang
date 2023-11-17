package main

import "fmt"

// NOTE ARRAYS ARE DIFFERENT THAN
// SLICES IN THAT THEY REQUIRE
// SIZE DURING THE DECLARATION ITSELF

func main() {
	fmt.Println("Welcome to Arrays!")
	arr := [4]string{"Hi", "hello", "hey"}
	fmt.Println(arr)
	fmt.Printf("len:%d", len(arr))
	fmt.Printf("\ncapacity:%d", cap(arr))
}
