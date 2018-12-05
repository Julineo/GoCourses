package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type task func(in, out chan interface{})

var agg uint32

func main() {
	tasks := []task{
		task(func(in, out chan interface{}) {
			fmt.Println("f1 start")
			out <- uint32(1)
			out <- uint32(2)
			out <- uint32(3)
			out <- uint32(100)
			fmt.Println("f1 end")
		}),
		task(func(in, out chan interface{}) {
			fmt.Println("f11 start")
			for val := range in {
				out <- val.(uint32) * 2
				time.Sleep(time.Millisecond * 250)
			}
			fmt.Println("f11 end")
		}),
		task(func(in, out chan interface{}) {
			fmt.Println("f2 start")
			for val := range in {
				out <- val.(uint32) * 3
				time.Sleep(time.Millisecond * 250)
			}
			fmt.Println("f2 end")
		}),
		task(func(in, out chan interface{}) {
			fmt.Println("f3 start")
			for val := range in {
				fmt.Println("f3 recieved: ", val)
				atomic.AddUint32(&agg, val.(uint32))
			}
			fmt.Println("f3 end")
		}),
	}

	Pipe(tasks...)

	time.Sleep(time.Millisecond * 50)

	fmt.Println("aggregated: ", agg)
}

func Pipe(fs ...task) {

	ch1 := make(chan interface{})
	wg := &sync.WaitGroup{}

	for _, f := range fs {
		ch2 := make(chan interface{})
		wg.Add(1)
		go func(f task, ch1, ch2 chan interface{}) {
			defer wg.Done()
			//defer close(ch2)
			f(ch1, ch2)
			close(ch2)
		}(f, ch1, ch2)
		ch1 = ch2
	}
	wg.Wait()
}


/*

package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go f(ch1, ch2)
	ch1 <- 5
	fmt.Println(<-ch2)
}

func f(ch1, ch2 chan int) {
	ch2 <- <-ch1
	close(ch2)
}

*/
