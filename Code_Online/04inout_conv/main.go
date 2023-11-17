package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Welcome to input conv")
	fmt.Println("Enter Your Rating")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	fmt.Println("Here is you Rating -> ", input)
	fmt.Printf("The Type of your rating is %T\n", input)

	if string_converted_to_int, err := strconv.ParseFloat(strings.TrimSpace(input), 64); err == nil {
		fmt.Println("here is one added to your input --> ", string_converted_to_int+1)
		fmt.Printf("\nAnd the type of your converted input now is %T", string_converted_to_int)
	}
}
