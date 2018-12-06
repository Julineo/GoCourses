package main

import (
	"fmt"
	"sync"
	"sort"
	"time"
)

func ExecutePipeline(fs ...job) {

	var ch1 chan interface{}
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
	fmt.Println("f1")

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for dat := range in {
		wg.Add(1)
		go singleHashHelper(dat, out, wg, mu)
	}

	wg.Wait()
}

func singleHashHelper(dat interface{}, out chan interface{}, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	data := fmt.Sprintf("%v", dat.(int))

	// mutex lock, so DataSignerMd5 run only one at a time
	mu.Lock()
	md5 := DataSignerMd5(data)
	mu.Unlock()

	tmpChan := make(chan string)
	// run DataSignerCrc32 concurently with another DataSignerCrc32
	go func(data string) {
		tmpChan <- DataSignerCrc32(data)
	}(data)

	cr32data := <-tmpChan

	out <- cr32data + "~" + DataSignerCrc32(md5)
}

func MultiHash(in, out chan interface{}) {
	fmt.Println("f2")
	ths := []string{"0", "1", "2", "3", "4", "5"}
	for val := range in {
		step := ""
		for _, th := range ths {
			tmp := DataSignerCrc32(th + val.(string))
			step += tmp
		}
		out <- step
	}
}

func CombineResults(in, out chan interface{}) {
	fmt.Println("f3")
	ar := []string{}
	for val := range in {
		tmp := val.(string)
		ar = append(ar, tmp)
	}
	sort.Strings(ar)
	fmt.Println("t1")
	time.Sleep(time.Millisecond * 5000)
	res := ar[0]
	fmt.Println("t2")
	for i := 1; i < len(ar); i++ {
		res += "_" + ar[i]
	}
	fmt.Println(res)
	out <- res
}
