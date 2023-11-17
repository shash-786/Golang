package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Welcome to the Slices sessions")
	var empty_arr []int
	fmt.Println(empty_arr, empty_arr == nil, len(empty_arr))

	new_arr := make([]int, 4)
	fmt.Println(new_arr, new_arr == nil, len(new_arr))

	str_arr := []string{"hello", "bye", "said"}

	str_arr = append(str_arr, "xqc", "joe rogan")
	fmt.Println(str_arr)

	twoD := make([][]int, 3)
	for i := 0; i < 3; i++ {
		twoD[i] = make([]int, 3)
		for j := 0; j < 3; j++ {
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			value, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
			if err != nil {
				fmt.Println(err)
				return
			}
			twoD[i][j] = int(value)
		}
	}

	fmt.Println(twoD)
}
