package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	fmt.Println("change to local start")
	b, err := ioutil.ReadFile("/etc/woda/GlobalConfigure_local.yaml")
	if err != nil {
		fmt.Print(err)
	}
	err = ioutil.WriteFile("/etc/woda/GlobalConfigure.yaml", b, 0644)
	if err != nil {
		fmt.Print("fail:")
		fmt.Print(err)
	} else {
		fmt.Println("success")
	}
	fmt.Println("change to local end")
}
