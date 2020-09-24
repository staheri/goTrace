package main

import (
	"fmt"
	"math/rand"
	"time"
)

// An Item is a stripped-down RSS item.
type Item struct{ Title, Channel, GUID string }

// A Fetcher fetches Items and returns the time when the next fetch should be
// attempted.  On failure, Fetch returns a non-nil error.
type Fetcher interface {
	Fetch() (items []Item, next time.Time, err error)
}

func Fetch(domain string) Fetcher {...} // fetches Items from domain

type Subscription interface {
    Updates() <-chan Item // stream of Items
    Close() error         // shuts down the stream
}

// sub implements the Subscription interface.
type sub struct {
    fetcher Fetcher   // fetches items
    updates chan Item // delivers items to the user
}

func Subscribe(fetcher Fetcher) Subscription {
  // Creates a new Subscription that repeatedly fetches items until Close is called.
    s := &sub{
        fetcher: fetcher,
        updates: make(chan Item), // for Updates
    }
    go s.loop()
    return s
}

func (s *sub) Updates() <-chan Item {
    return s.updates
}

func (s *sub) Close() error {
    // TODO: make loop exit
    // TODO: find out about any error
    return err
}

// loop fetches items using s.fetcher and sends them
// on s.updates.  loop exits when s.Close is called.
func (s *sub) loop() {...}

func Merge(subs ...Subscription) Subscription {...} // merges several streams


func main() {
    // Subscribe to some feeds, and create a merged update stream.
    merged := Merge(
        Subscribe(Fetch("blog.golang.org")),
        Subscribe(Fetch("googleblog.blogspot.com")),
        Subscribe(Fetch("googledevelopers.blogspot.com")))

    // Close the subscriptions after some time.
    time.AfterFunc(3*time.Second, func() {
        fmt.Println("closed:", merged.Close())
    })

    // Print the stream.
    for it := range merged.Updates() {
        fmt.Println(it.Channel, it.Title)
    }

    panic("show me the stacks")
}
