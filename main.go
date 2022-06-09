package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	boundedWorkPooling()
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

// drop pattern important pattern for service that may experience
// loads of trafic at times, letting you drop request when capacity
// is full. The key of this pattern is the default statement that
// lets you drop, cancel or redirect work when capacity is full.
// By that it give you options
func drop() {
	cap := 4
	ch := make(chan string, cap)

	go func() {
		for p := range ch {
			fmt.Println("child: recv signal :", p)
		}
	}()

	const work = 2000
	for w := 0; w < work; w++ {
		select {
		case ch <- "data":
			fmt.Println("Parent: sent signal:", w)
		default:
			fmt.Println("parent drop data:", w)
		}
	}

	close(ch)
	fmt.Println("parent: sent shutdown signal")
	time.Sleep(time.Second)
	fmt.Println("--------------------------------------------")
}

// cancelation pattern is used to tell a function performing some i/o
// how long is willing to whait by canceling it or just walking away
// In this example its important to create a buffered chanel of 1
// because if we walk away the goroutine needs to be able to have
// someone to recieve its works or else it blocks for ever and becomes
// a memory leak
func cancelation() {
	duration := 150 * time.Microsecond
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	defer cancel()

	ch := make(chan string, 1)

	go func() {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		ch <- "data"
	}()

	select {
	case d := <-ch:
		fmt.Println("work complete", d)
	case <-ctx.Done():
		fmt.Println("work canceled")
	}
	time.Sleep(time.Second)
	fmt.Println("--------------------------------------------")
}

// Fan In/Out Semaphore pattern provides a mechanics to control
// the number of goroutine executing work at any given time while
// still creating a unique goroutine for each piece of works
func fanOutSemaphore() {
	children := 2000
	ch := make(chan string, children)

	g := runtime.GOMAXPROCS(0)
	sem := make(chan bool, g)

	for c := 0; c < children; c++ {
		go func(child int) {
			sem <- true
			{
				t := time.Duration(rand.Intn(200)) * time.Microsecond
				time.Sleep(t)
				ch <- "data"
				fmt.Println("child : sent signal :", child)
			}
			<-sem
		}(c)
	}

	for children > 0 {
		d := <-ch
		children--
		fmt.Println(d)
		fmt.Println("parent: recv signal :", children)
	}

	time.Sleep(time.Second)
	fmt.Println("--------------------------------------------")
}

// BoundedWorkPooling is a pattern that uses pool of goroutines to
// perform a fixed amount of know work
func boundedWorkPooling() {
	work := []string{"paper", "paper", "paper", "paper", 2000: "paper"}

	g := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan string, g)

	for c := 0; c < g; c++ {
		go func(child int) {
			defer wg.Done()
			for wrk := range ch {
				fmt.Printf("child %d: recvd signal :%s\n", child, wrk)
			}
			fmt.Printf("child %d: recvd shutdown signal\n", child)
		}(c)
	}

	for _, wrk := range work {
		ch <- wrk
	}
	close(ch)
	wg.Wait()

	time.Sleep(time.Second)
	fmt.Println("--------------------------------------------")
}
