package main

import (
	"bytes"
	"fmt"
	"math/rand"
)

func main() {

	linkls := make([]bytes.Buffer, 0)
	for i := 0; i < 100; i++ {
		var tmp bytes.Buffer
		tmp.WriteString(getRandomStr(4))
		linkls = append(linkls, tmp)
	}
	//数组
	linkmap := map[int64]string{}

	for k,v:= range linkls{
		fmt.Println(v.String())
	}

}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func getRandomStr(n int64) string {

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)

}
