package main

import "fmt"

func main() {
	fmt.Println("Welcome to Range Session")

	var sum int = 0
	arr := []int{1, 2, 3}

	for index, value := range arr {
		sum += value
		fmt.Printf("%d : %d\n", index, value)
	}
	fmt.Println("sum = ", sum)
}
