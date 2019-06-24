package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("------------------------")
	cnum := 100000000
	//并行
	ch := make(chan []int, 4)
	fmt.Println("并行 for start")
	t1 := time.Now().UnixNano()
	//分4个线程
	count := cnum / 4
	for i := 0; i < 4; i++ {
		go func(i int) {
			start := i * count
			end := start + count
			ils := make([]int, 0)
			for j := start; j < end; j++ {
				//fmt.Println("i:" + strconv.Itoa(i) + "," + "j:" + strconv.Itoa(j))
				ils = append(ils, j)
			}
			ch <- ils
		}(i)
	}
	ls1 := <-ch
	ls2 := <-ch
	ls3 := <-ch
	ls4 := <-ch
	ls := make([][]int, 0)
	ls = append(ls, ls1)
	ls = append(ls, ls2)
	ls = append(ls, ls3)
	ls = append(ls, ls4)
	t2 := time.Now().UnixNano()
	fmt.Println("组织数据 time ms:", (t2-t1)/1000000)
	chout := make(chan int, 4)
	for i := 0; i < 4; i++ {
		go func(i int) {
			out := 0
			item := ls[i]
			for _, v := range item {
				out += v
			}
			//fmt.Println(out)
			chout <- out
		}(i)
		//time.Sleep(3*time.Second)
	}
	one := <-chout
	two := <-chout
	thr := <-chout
	fur := <-chout
	tol := one + two + thr + fur
	t3 := time.Now().UnixNano()
	fmt.Println("数据相加 time ms:", (t3-t2)/1000000)
	fmt.Println("sum(i):", tol)
	fmt.Println("并行 for end")
	fmt.Println("------------------------")
}
