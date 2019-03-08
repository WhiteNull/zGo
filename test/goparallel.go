package main

import (
	"fmt"
	"time"
)

func main() {
	cnum := 100000000
	//并行
	ch := make(chan []int, 4)
	ils := make([]int, 0)
	fmt.Println("并行 for start")
	t1 := time.Now().UnixNano()
	//分4个线程
	count := cnum / 4
	for i := 0; i < 4; i++ {
		go func(i int) {
			start := i*count + 1
			end := start + count + 1
			for j := start; j < end; j++ {
				ils = append(ils, j)
			}
			ch <- ils
		}(i)
	}
	for {
		select {
		case <-ch:
			fmt.Println(len(<-ch))
		default:
		}
	}
	t2 := time.Now().UnixNano()
	fmt.Println("time ms:", (t2-t1)/1000000)
	//fmt.Println("sum(i):", tol)
	fmt.Println("并行 for end")

}
