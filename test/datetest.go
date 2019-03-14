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
	enddt := "2019-02-21"
	tEndDt, _ := time.Parse("2006-01-02", enddt)
	if tEntryDt.After(tEndDt) {
		fmt.Println("enter")
	} else {
		fmt.Println("not enter")
	}

}
