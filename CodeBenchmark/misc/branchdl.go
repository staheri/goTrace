package main

func S(ch chan int, done chan struct{}) {
	ch <- 1
	done <- struct{}{}
}

func R(ch chan int, done chan struct{}) {
	<-ch
	done <- struct{}{}
}

func main() {
	done := make(chan struct{})
  ch := make(chan int)
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			go S(ch, done)
		} else {
			go R(ch, done)
		}
	}
	<-done
}
