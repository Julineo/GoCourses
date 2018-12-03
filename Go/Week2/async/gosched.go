package main

import (
	"fmt"
	"runtime"
	//"runtime"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		runtime.Gosched() //try to comment
		fmt.Println(s)
	}
}

func main() {
	go say("world")
	say("hello")
	//	fmt.Scanln() //try to comment
}
