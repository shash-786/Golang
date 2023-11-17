package main

import "fmt"

func find(num int, nums ...int) {
	for index, value := range nums {
		if value == num {
			fmt.Println("Found ", num, " at ", index, "in", nums)
			return
		}
	}

	fmt.Printf("Not found %d in %v\n", num, nums)
}

func main() {
	find(1, 1, 2, 3, 4)
	find(3, 1, 2, 2, 2, 2, 4, 3)
	find(2, 1)
	find(1)
}
