// Example program to show a goroutine leak. We start
// goroutines that range over a channel, but nothing ever
// closes that channel.
package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func main() {

	// Capture starting number of goroutines.
	startingGs := runtime.NumGoroutine()

	names := []string{"Anna", "Jacob", "Kell", "Carter", "Rory"}
	processRecords(names)

	// Hold the program from terminating for 1 second to see
	// if any goroutines created by processRecords terminate.
	time.Sleep(time.Second)

	// Capture ending number of goroutines.
	endingGs := runtime.NumGoroutine()

	// Report the results.
	fmt.Println("========================================")
	fmt.Println("Number of goroutines before:", startingGs)
	fmt.Println("Number of goroutines after :", endingGs)
	fmt.Println("Number of goroutines leaked:", endingGs-startingGs)
}

// processRecords is given a slice of values such as lines
// from a file. The order of these values is not important
// so the function can start multiple workers to perform some
// processing on each record then feed the results back.
func processRecords(records []string) {

	// Load all of the records into the input channel. It is
	// buffered with just enough capacity to hold all of the
	// records so it will not block.

	total := len(records)
	input := make(chan string, total)
	for _, record := range records {
		input <- record
	}
	close(input) // What if we forget to close the channel?

	// Start a pool of workers to process input and send
	// results to output. Base the size of the worker pool on
	// the number of logical CPUs available.

	output := make(chan string, total)
	workers := runtime.NumCPU()
	for i := 0; i < workers; i++ {
		go worker(i, input, output)
	}

	// Receive from output the expected number of times. If 10
	// records went in then 10 will come out.

	for i := 0; i < total; i++ {
		result := <-output
		fmt.Printf("[result  ]: output %s\n", result)
	}
}

// worker is the work the program wants to do concurrently.
// This is a blog post so all the workers do is capitalize a
// string but imagine they are doing something important.
//
// Each goroutine can't know how many records it will get so
// it must use the range keyword to receive in a loop.
func worker(id int, input <-chan string, output chan<- string) {
	for v := range input {
		fmt.Printf("[worker %d]: input %s\n", id, v)
		output <- strings.ToUpper(v)
	}
	fmt.Printf("[worker %d]: shutting down\n", id)
}

