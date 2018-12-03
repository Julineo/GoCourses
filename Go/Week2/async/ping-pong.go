package main

import (
	"fmt"
	"time"
)

func main() {
	var ball int
	table := make(chan int)
	go player(table)
	go player(table)

	table <- ball
	time.Sleep(1 * time.Second)
	<-table
}

func player(table chan int) {
	for {
		ball := <-table
		ball++
		fmt.Println(ball)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}
