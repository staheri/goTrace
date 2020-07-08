package main

import (
	"fmt"
	"io/ioutil"
	"os"
	)


// function declaration
func fibonacci(c, quit chan int, f interface{}) {
	x, y := 0, 1
	for {
		// select statement
		select {
		// case statement
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			//return statment
			return
		}
	}
}

func main() {

	f, err := os.Create("sample-trace-t0.txt")
	check(err)
	defer f.Close()
	var trace string
	trace = ""

	//channel creation
	c := make(chan int)
	trace = trace + fmt.Sprintf("main>New Channel\n\tName: %s\n\tType: %s\n\tCapacity: %s\n\tLocation: %s\n","c","int","default","fibonacci.go:35")

	quit := make(chan int)
	trace = trace + fmt.Sprintf("main>New Channel\n\tName: %s\n\tType: %s\n\tCapacity: %s\n\tLocation: %s\n","quit","int","default","fibonacci.go:38")
	// goroutine spawn
	trace = trace + fmt.Sprintf("main>Go func1\n\tLocation: %s\n","fibonacci.go:43")
	go func() {
		// for statment
		for i := 0; i < 10; i++ {
			// print + channel receive
			trace = trace + fmt.Sprintf("func1 > Recv Channel\n\tName: %s\n\tType: %s\n\tCapacity: %s\n\tLocation: %s\n","quit","int","default","fibonacci.go:")
			fmt.Println(<-c)
		}
		//channel send
		quit <- 0
	}()
	fibonacci(c, quit)
}
