package main

//import "fmt"
import (
  "os"
  "runtime/trace"
  "log"
)

func main() {
  f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()


  ctx, task := trace.NewTask(ctx, "makeCappuccino")
  trace.Log(ctx, "orderID", "0")

  milk := make(chan bool)
  espresso := make(chan bool)

  go func() {
          trace.WithRegion(ctx, "steamMilk", "1")
          milk <- true
  }()
  go func() {
          trace.WithRegion(ctx, "extractCoffee", "2")
          espresso <- true
  }()
  go func() {
          defer task.End() // When assemble is done, the order is complete.
          <-espresso
          <-milk
          trace.WithRegion(ctx, "mixMilkCoffee", "3")
  }()
}
