package main

import "fmt"

func main() {
	aa := "123456789012345678"
	fmt.Println("****" + aa[4:14] + "****")

	bb := "1234567890123"
	fmt.Println(bb[0:3] + "****" + bb[len(bb)-4:])

	cc := "2018-08-06 07:14:31"
	fmt.Println(len(cc))
	fmt.Println(cc[0:19])
}
