package main

import (
	"crypto/md5"
	"fmt"
	log "github.com/xiaomi-tc/log15"
	"strconv"
	"time"
)

func main() {

	md5str := strconv.FormatInt(327955, 10) + "2018-11" + "SelectPayroll"
	log.Debug("payrollH5Select", "md5 之前的值:", md5str)
	str := md5.Sum([]byte(md5str))
	md5value := fmt.Sprintf("%x", str)
	log.Debug("payrollH5Select", "md5 之后的值:", md5value)

	fmt.Println("--------------------------")

	cnum := 100000000
	//串行
	ils := make([]int, 0)
	fmt.Println("串行 for start")
	t1 := time.Now().UnixNano()
	for i := 0; i < cnum; i++ {
		ils = append(ils, i)
	}
	tol := 0
	for i := 0; i < len(ils); i++ {
		tol += i
	}
	t2 := time.Now().UnixNano()
	fmt.Println("time ms:", (t2-t1)/1000000)
	fmt.Println("sum(i):", tol)
	fmt.Println("串行 for end")
	fmt.Println("------------------------")

	//并行
	ch := make(chan []int, 4)
	ils = make([]int, 0)
	fmt.Println("并行 for start")
	t1 = time.Now().UnixNano()
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
	//chils := <-ch
	//tol = 0
	//for i := 0; i < len(chils); i++ {
	//	tol += i
	//}
	t2 = time.Now().UnixNano()
	fmt.Println("time ms:", (t2-t1)/1000000)
	fmt.Println("sum(i):", tol)
	fmt.Println("并行 for end")
}
