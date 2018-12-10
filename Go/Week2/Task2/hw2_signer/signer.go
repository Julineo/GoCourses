package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

const th = 6

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

	wg := &sync.WaitGroup{}
	mux := &sync.Mutex{}

	for dat := range in {
		wg.Add(1)
		go singleHashWorker(dat, out, wg, mux)
	}

	wg.Wait()
}

func singleHashWorker(dat interface{}, out chan interface{}, wg *sync.WaitGroup, mux *sync.Mutex) {
	data := fmt.Sprintf("%v", dat.(int))
	start := time.Now()

	mux.Lock()
	md5Data := DataSignerMd5(data)
	mux.Unlock()
	// crc32md5Data concurently with variable (not correct) you'll  get a race
	/*	crc32md5Data := ""
		go func() {
			crc32md5Data = DataSignerCrc32(md5Data)
		}()
	*/
	// crc32md5Data concurently with channel (correct implementation)
	tmpChan := make(chan string)
	go func() {
		tmpChan <- DataSignerCrc32(md5Data)
	}()

	start = time.Now()

	crc32Data := DataSignerCrc32(data)

	//	fmt.Println("data: ", data, "crc32Data: ", crc32Data)
	//	fmt.Println("data: ", data, "md5Data: ", md5Data)
	crc32md5Data := <-tmpChan
	out <- crc32Data + "~" + crc32md5Data
	wg.Done()
	fmt.Println(data, time.Since(start))
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for dat := range in {
		wg.Add(1)
		go multiHashWorker(dat, out, wg)
	}

	wg.Wait()
}

func multiHashWorker(dat interface{}, out chan interface{}, wg *sync.WaitGroup) {
	start := time.Now()
	data := dat.(string)
	//	mux := &sync.Mutex{}
	wg2 := &sync.WaitGroup{}

	arr := make([]string, th)
	wg2.Add(th)
	for i := 0; i < th; i++ {
		i := i
		go func() {
			icrc32 := DataSignerCrc32(fmt.Sprintf("%v", i) + data)

			//			mux.Lock()
			arr[i] = icrc32
			//	fmt.Println("data: ", data, "th: ", i, "ss: ", arr[i])
			//			mux.Unlock()

			wg2.Done()
		}()
	}
	wg2.Wait()
	res := strings.Join(arr, "")

	out <- res
	wg.Done()
	fmt.Println("multi: ", time.Since(start))
}

func CombineResults(in, out chan interface{}) {
	fmt.Println("f3")
	ar := []string{}
	for val := range in {
		tmp := val.(string)
		ar = append(ar, tmp)
	}
	sort.Strings(ar)
	res := ar[0]
	for i := 1; i < len(ar); i++ {
		res += "_" + ar[i]
	}
	out <- res
}
