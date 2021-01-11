package main

// https://github.com/moby/moby/pull/25384
// buggy version:  87e48ecd048c0b083fe09fb8d74c83364abd41e6
// On each itertion, there is a wait causing deadlock
// fix: Put wait outside of loop

import (
  "time"
  "sync"
  "fmt"
)

func main() {
    var group sync.WaitGroup
    var t = []int{1, 2, 3, 4}
    group.Add(len(t))
    for _,p := range t{
      go func(p int){
        fmt.Println(p)
        defer group.Done()
      }(p)
      group.Wait()
    }
    // fix: group.Wait()

    time.Sleep(time.Millisecond*10)
    fmt.Println("End of main!")
}
