package main

import "fmt"

func fun(num int) {
	num = num + 1
}

func fun2(num *int) {
	*num = *num + 1
}

func main() {
	fmt.Println("Welcome To Session on Pointers")
	var num int = 4

	fun(num)
	fmt.Println("Num ->", num)
	fun2(&num)
	fmt.Println("Num ->", num)
}
