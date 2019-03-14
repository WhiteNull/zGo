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

}
