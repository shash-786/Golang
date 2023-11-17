package main

import "fmt"

func fact(num int) int {
	if num == 0 || num == 1 {
		return 1
	}
	return num * fact(num-1)
}

func main() {
	for i := 0; i < 9; i++ {
		fmt.Printf("fact of %d is %d\n", i, fact(i))
	}
}
