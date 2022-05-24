package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	pooling()
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

//  fanOut pattern that creates a Goroutine for each piece of work
//  that is pending and can be done concurrently
func fanOut() {
	children := 2000
	ch := make(chan string, children)

	for c := 0; c < children; c++ {
		go func(child int) {
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			ch <- "data"

			fmt.Println("child: sent signal:", child)
		}(c)
	}

	for children > 0 {
		d := <-ch
		children--

		fmt.Println(d)
		fmt.Println("parent: recieved signal:", children)
	}
}

// waitForTask is a foundational patternused by larger patterns like pooling
// it creates an unbuffered chanel so there is a guaranteee at the
// signal level, this is critical importante so its posible to add
// mechanics later like canceling or timeouts. Once the chanel is,
// a child Goroutine is created whaiting for the signal to perform
// work send by the parent.
func waitForTask() {
	ch := make(chan string)

	go func() {
		d := <-ch
		fmt.Println("child:recieved signal:", d)
	}()

	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	ch <- "data"

	fmt.Println("parent:send signal")
	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

// pooling pattern uses the waitForTask pattern. This pattern allows me to
// manage resource usage across a well defined number of Gourutines.In Go
// pooling is not needed for efficiency in CPU processing like at the operating
// system. Its more important for efficiency in resource usage.
func pooling() {
	ch := make(chan string)

	g := runtime.GOMAXPROCS(0)
	for c := 0; c < g; c++ {
		go func(child int) {
			for d := range ch {
				fmt.Printf("child %d: recv'd signal: %s\n", child, d)
			}

			fmt.Printf("child %d, recv'd shutdown signal\n", child)
		}(c)
	}

	const work = 100
	for w := 0; w < work; w++ {
		ch <- "data"
		fmt.Println("parent:sent signal:", w)
	}

	close(ch)
	fmt.Println("parent:sent shutdown signal")
	time.Sleep(time.Second)
	fmt.Println("--------------------------")
}
