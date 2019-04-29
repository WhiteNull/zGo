package main

import (
	"fmt"
	"time"
)

func main() {
	str := "0000-00-00"
	tEntryDt, err := time.Parse("2006-01-02", str)
	if err != nil {
		fmt.Println("time parse err")
	}
	enddt := "2019-03-23"
	tEndDt, _ := time.Parse("2006-01-02", enddt)
	if tEntryDt.After(tEndDt) {
		fmt.Println("enter")
	} else {
		fmt.Println("not enter")
	}

	aa := "2019-03-23"
	bb := "2019-03-24 12:32:34"
	taa, _ := time.Parse("2006-01-02", aa)
	tbb, _ := time.Parse("2006-01-02 15:04:05", bb)
	if taa.After(tbb) {
		fmt.Println("aaaaaa")
	} else {
		fmt.Println("bbbbbbb")
	}
	fmt.Println(taa)
	fmt.Println(tbb)

}
