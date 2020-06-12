package main

import "net"
import "fmt"

func handler(c net.Conn) {
    c.Write([]byte("ok"))
    fmt.Println("ok")
    c.Close()
}

func main() {
    l, err := net.Listen("unix", "test")
    if err != nil {
        panic(err)
    }
    for {
        c, err := l.Accept()
        if err != nil {
            continue
        }
	fmt.Println("for")
        go handler(c)
    }
}

