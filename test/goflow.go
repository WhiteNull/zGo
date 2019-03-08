package main

import (
	"fmt"
	"time"
)

func main() {
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
}
