package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	welcome := "Welcome to Input Practical"
	fmt.Println(welcome)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter You rating!")
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println("your Rating Was --> ", input)
	fmt.Printf("the Datatype of input was %T", input)
}
