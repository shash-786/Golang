package main

import (
	"fmt"
	"log"
	"reflect"
)

func main() {
	m1 := make(map[string]int32)
	m2 := make(map[int]float64)

	m1["a"] = 1
	m1["b"] = 2
	m1["c"] = 3

	m2[1] = 1
	m2[2] = 2.5
	m2[3] = 3

	print("Sum:", sum(m1))
	print("Sum:", sum(m2))
	print("Sum:", []string{"apple"})
	print("Sum:", nil)
}

func sum[K comparable, V int32 | float64](m map[K]V) V {
	var sum V
	for _, v := range m {
		sum += v
	}
	return sum
}

func print(prefix string, s any) {
	t := reflect.TypeOf(s)
	if t == nil {
		log.Fatalf("Recieved Nil!\n")
	}

	switch t.Kind() {
	case reflect.Int64:
	case reflect.Int32:
		fmt.Printf(prefix+" %d\n", s)
	case reflect.Float64:
		fmt.Printf(prefix+" %.2f\n", s)
	default:
		fmt.Printf("recieved unknown type: %v\n", t)
	}
}
