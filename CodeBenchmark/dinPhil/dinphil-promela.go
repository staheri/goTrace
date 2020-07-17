// +build ignore

package main

import (
	"fmt"
	"time"
)

const (
	NO           = 1
	YES          = 2
	ARE_YOU_FREE = 3
	RELEASE      = 4
)

func checkVal(val,exp int,){
	if val != exp{
		fmt.Printf("Expected %v got %v\n",exp,val)
		panic("Invalid value")
	}
}

func Fork(lch,rch chan int) {
	for {
		select{
		case rval := <- rch:
			checkVal(rval,ARE_YOU_FREE)
			rch <- YES
			for{
				select{
				case lval := <- lch:
					if lval != ARE_YOU_FREE{
						lch <- lval
					}else{
						lch <- NO
					}
					//checkVal(lval,ARE_YOU_FREE)

				case rval := <- rch:
					//checkVal(rval,RELEASE)
					if rval != RELEASE{
						//fmt.Printf("CHAN:WRONG VALUE RECEIVED. RESENDING\n")
						rch <- rval
					}else{
						break
						}
				}
			}
		case lval := <- lch:
			checkVal(lval,ARE_YOU_FREE)
			lch <- YES
			for{
				select{
				case rval := <- rch:
					//checkVal(rval,ARE_YOU_FREE)
					if rval != ARE_YOU_FREE{
						//fmt.Printf(">CHAN:WRONG VALUE RECEIVED. RESENDING\n")
						rch <- rval
					}else{
						rch <- NO
					}

				case lval := <- lch:
					if lval != RELEASE{
						//fmt.Printf(">>CHAN:WRONG VALUE RECEIVED. RESENDING\n")
						lch <- lval
					}else{
						break
					}
					//checkVal(lval,RELEASE)
					//break
				}
			}
		}
	}
}

func phil(lch, rch chan int, id int) {
	for {
		for {
			//fmt.Printf("Phil %v asking if left fork is free\n",id)
			lch <- ARE_YOU_FREE
			//fmt.Printf("Phil %v waiting to hear from left fork\n",id)
			lval := <- lch
			if lval == YES{
				//fmt.Printf("Phil %v left fork free...Break! now check right fork\n",id)
				break
			} else{
				//fmt.Printf("Phil %v left fork busy\n",id)
				continue
			}
		}
		for {
			//fmt.Printf("Phil %v asking if right fork is free\n",id)
			rch <- ARE_YOU_FREE
			//fmt.Printf("Phil %v waiting to hear from right fork\n",id)
			rval := <- rch
			if rval == YES{
				fmt.Printf("\t\t\t\t\t\t\t\tPhil %v right fork free\nPhil %v is eating\n",id,id)
				//panic("SUCCESS")
				time.Sleep(2*time.Nanosecond)
				//fmt.Printf("Phil %v releasing left fork\n",id)
				lch <- RELEASE
				//fmt.Printf("Phil %v releasing right fork\n",id)
				rch <- RELEASE
				//fmt.Printf("Phil %v both forks released\n",id)
				break
			} else if rval == NO{
				//fmt.Printf("Phil %v right fork busy...Releasing left fork\n",id)
				lch <- RELEASE
				//fmt.Printf("Phil %v left fork released\n",id)
				break
			} else{
				panic("WEIRD")
			}
		}
	}
}

func main() {
	ch0 := make(chan int)
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	ch4 := make(chan int)
	ch5 := make(chan int)
	go phil(ch5, ch0, 0)
	go phil(ch1, ch2, 1)
	go phil(ch3, ch4, 2)
	go Fork(ch0, ch1)
	go Fork(ch2, ch3)
	go Fork(ch4, ch5)
	time.Sleep(1*time.Microsecond)
}
