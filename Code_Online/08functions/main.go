package main

import "fmt"

func add(a, b int) int {
	return a + b
}

func main() {
	num1 := 3
	num2 := 5

	fmt.Println(add(num1, num2))
}
