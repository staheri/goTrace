package main

var count = 0

func main() {
	go incrementCount()
	go incrementCount()
}

func incrementCount() {
	if count == 0{
		count++
	}
}
