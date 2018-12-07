package main

import (
	"fmt"
	"sort"
	"sync"
	"strings"
	"strconv"
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
	crc32md5Data := <-tmpChan

	crc32Data := DataSignerCrc32(data)

	//	fmt.Println("data: ", data, "crc32Data: ", crc32Data)
	//	fmt.Println("data: ", data, "md5Data: ", md5Data)
	out <- crc32Data + "~" + crc32md5Data
	wg.Done()
}

func MultiHash(in, out chan interface{}) {

	wg := &sync.WaitGroup{}

	for dat := range in {
		wg.Add(1)
		go multiHashWorker(dat, out, wg)
	}

	wg.Wait()
}

/*
func multiHashWorker(dat interface{}, out chan interface{}, wg *sync.WaitGroup) {
	data := dat.(string)
	mux := &sync.Mutex{}
	wg2 := &sync.WaitGroup{}

	ths := []string{"0", "1", "2", "3", "4", "5"}
	step := ""
	for _, th := range ths {
		wg2.Add(1)
		th := th
		go func() {
			crc32 := DataSignerCrc32(th + data)

			mux.Lock()
			step += crc32
			fmt.Println("data: ", data, "th: ", th, "ss: ", step)
			mux.Unlock()

			wg2.Done()
		}()
	}
	wg2.Wait()
	out <- step
	wg.Done()
}
*/

func multiHashWorker(in interface{}, out chan interface{}, wg *sync.WaitGroup) {
	Th := 6
	defer wg.Done()
	mu := &sync.Mutex{}
	wgCrc32 := &sync.WaitGroup{}

	concatArray := make([]string, Th)
//	step := ""
	for i := 0; i < Th; i++ {
		wgCrc32.Add(1)
		data := strconv.Itoa(i) + in.(string)
		go func(data string, index int, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			data = DataSignerCrc32(data)
			mu.Lock()
			concatArray[index] = data
			//step += data
			fmt.Printf("%s MultiHash: crc32(th+step1)) %d %s\n", in, index, data)
			mu.Unlock()
		}(data, i, wgCrc32, mu)
	}
	wgCrc32.Wait()
	result := strings.Join(concatArray, "")
//	result := step
	fmt.Printf("%s MultiHash result: %s\n", in, result)
	out <- result
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
