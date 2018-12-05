package main

import (
	"fmt"
	"sync"
)

func ExecutePipeline(fs ...job) {

	ch1 := make(chan interface{})
	wg := &sync.WaitGroup{}

	for _, f := range fs {

		ch2 := make(chan interface{})
		wg.Add(1)

		go func(f job, ch1, ch2 chan interface{}) {
			defer wg.Done()
			f(ch1, ch2)
			close(ch2)
		}(f, ch1, ch2)

		ch1 = ch2
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	fmt.Println(1)
	out <- 1
}

func MultiHash(in, out chan interface{}) {
	fmt.Println(2)
	out <- 2
}

func CombineResults(in, out chan interface{}) {
	fmt.Println(3)
	out <- 3
}
