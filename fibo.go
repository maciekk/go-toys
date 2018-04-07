package main

import "fmt"

func fibonacci(ch chan int) {
	x, y := 0, 1
	for {
		ch <- x
		x, y = y, x+y
	}
}

func main() {
	ch := make(chan int)
	go fibonacci(ch)
	var x int
	for i := 0; i < 20; i++ {
		x = <- ch
		fmt.Printf("%v\n", x)
	}
}
