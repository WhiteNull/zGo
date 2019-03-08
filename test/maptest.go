package main

import (
	"crypto/md5"
	"fmt"
	log "github.com/xiaomi-tc/log15"
	"strconv"
)

func main() {

	md5str := strconv.FormatInt(327955, 10) + "2018-11" + "SelectPayroll"
	log.Debug("payrollH5Select", "md5 之前的值:", md5str)
	str := md5.Sum([]byte(md5str))
	md5value := fmt.Sprintf("%x", str)
	log.Debug("payrollH5Select", "md5 之后的值:", md5value)

	fmt.Println("--------------------------")

	aa := map[string]string{}
	aa["a"] = "a"
	aa["b"] = "b"
	aa["c"] = "c"

	fmt.Println(len(aa))

	index := 0
	for k, v := range aa {
		index++
		fmt.Println(k)
		fmt.Println(v)
		if index == len(aa) {
			fmt.Println("11111111")
		}

	}

	fmt.Println("before delete")
	fmt.Println(aa)
	for k := range aa {
		if k == "b" {
			delete(aa, "b")
		}
		fmt.Println("kkkkkkk")
		fmt.Println(k)
	}
	fmt.Println("after delete")
	fmt.Println(aa)

}
