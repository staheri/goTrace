package main

import (
	"fmt"
	"sync"
	"runtime"
	//"time"
)

func main(){
	runtime.GOMAXPROCS(2)
	var wg sync.WaitGroup
	for i:=0;i<5;i++{
		wg.Add(1)
		go func(ii int){
			//time.Sleep(11*time.Millisecond)
			runtime.Gosched()
			fmt.Println(ii)

			wg.Done()
		}(i)
	}
	wg.Wait()
}
