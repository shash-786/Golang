package main

import (
	"fmt"
	"io"
	"log"
)

type myreaeder struct {
  str string
  pos int
}

func (mr *myreaeder) Read(p []byte) (n int, err error) {
  if mr.pos == len(mr.str) {
    return 0, io.EOF
  }
  n = copy(p, mr.str[mr.pos:mr.pos+1])
  mr.pos++;
  return n, nil
}

func main() {
  mr := &myreaeder{
    str: "Hello World",
    pos: 0,
  }

  output, err := io.ReadAll(mr)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(string(output))
}
