package main

// сюда писать код

import (
	"fmt"
)

func ExecutePipeline(fs ...job) {

	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	for _, f := range fs {
		ch1, ch2 = ch2, ch1
		go f(ch1, ch2)
	//	<-ch2
	}
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
