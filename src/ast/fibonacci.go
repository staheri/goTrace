package main

import "fmt"

// function declaration
func fibonacci(c, quit chan int) {

	x, y := 0, 1
	LOOP:
		for {
			// select statement
			select {
				// case statement
			case c <- x:
				x, y = y, x+y
			case <-quit:
				fmt.Println("quit")
				//return statment
				//return
				break LOOP
		}
	}
	return
}

func main() {
	//channel creation
	c := make(chan int)
	quit := make(chan int)
	// goroutine spawn
	go func() {
		// for statment
		for i := 0; i < 10; i++ {
			// print + channel receive
			fmt.Println(<-c)
		}
		//channel send
		quit <- 0
	}()
	fibonacci(c, quit)
}
