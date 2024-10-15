package main

import (
	"fmt"
	"io"
	"log"
)

type myreader struct {
	str string
	pos int
}

func (mr *myreader) Read(p []byte) (n int, err error) {
	if mr.pos == len(mr.str) {
		return 0, io.EOF
	}
	n = copy(p, mr.str[mr.pos:mr.pos+1])
	mr.pos++
	return n, nil
}

func main() {
	mr := &myreader{
		str: "Hello World",
		pos: 0,
	}

	output, err := io.ReadAll(mr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
}
