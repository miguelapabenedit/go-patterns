package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	waitForResult()
}

// A Goroutine foundation pattern used by larger patterns like fan out/in.In this patttern a gourtine is
// created to performe some know work and signals their result back to the gourtine that created them.
// This allows for the actual work to be placed  on a goroutine that can be terminated or walked away from
func waitForResult() {
	ch := make(chan string)

	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Microsecond)
		ch <- "data"

		fmt.Println("child: sent signal")
	}()

	d := <-ch
	fmt.Println("parent: recv`d signal:", d)
	time.Sleep(time.Second)
	fmt.Println("--------------------------------------------")
}
