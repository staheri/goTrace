package main

//import "fmt"

func main() {
	c := make(chan struct{}) // or buffered channel
	go func() { c <- struct{}{} }()
	close(c)
}
