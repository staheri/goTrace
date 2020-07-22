// +build ignore

package main

import (
	_"fmt"
	"time"
)

const N = 3

type Philosopher struct {
	id                int
	thinkingTime      int // Millisecond
	eatingTime        int // Millisecond
}

func Fork(fork *int, ch chan int) {
	for {
		*fork = 2
		ch <- 1
		<-ch
	}
}

func phil(fork1, fork2 *int, ch1, ch2 chan int, philStruct *Philosopher) {
	for {
		time.Sleep(time.Duration(philStruct.thinkingTime)*time.Millisecond)
		select {
		case <-ch1:
			select {
			case <-ch2:
				//fmt.Printf("phil %d got both fork\n", philStruct.id)
				time.Sleep(time.Duration(philStruct.eatingTime)*time.Millisecond)
				ch1 <- *fork1
				ch2 <- *fork2
			default:
				ch1 <- *fork1
			}
		case <-ch2:
			select {
			case <-ch1:
				//fmt.Printf("phil %d got both fork\n", philStruct.id)
				time.Sleep(time.Duration(philStruct.eatingTime)*time.Millisecond)
				ch2 <- *fork2
				ch1 <- *fork1
			default:
				ch2 <- *fork2
			}
		}
	}
}

func main() {
	// Initilizing Forks
	var forks [N]int
	for i := range forks {
  	forks[i] = 1
	}

	// Initializing Channels
	var chans [N]chan int
	for i := range chans {
  	chans[i] = make(chan int)
	}

	// Initilizing Philosophers
	var phils = [N]Philosopher{
		0:  {0,1,1},
	  1:  {1,2,2},
		2:  {2,3,3},
	}

	for i:=0 ; i<N ; i++{
		go phil(&forks[i],&forks[(i+1)%N],chans[i],chans[(i+1)%N],&phils[i])
	}
	for i:=0 ; i<N ; i++{
		go Fork(&forks[i],chans[i])
	}

	time.Sleep(100*time.Millisecond)
}
