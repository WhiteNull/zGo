package main

import (
	"fmt"
	"math"
)

var A  = "12345678"
func main() {

	A = "bb"
	fmt.Println(A)
	pa()


}

func pa (){
	fmt.Println(A)
}

//小数向下取整保留2位小数
//f 小数
func FloatDown2(f float64) float64 {
	return math.Floor(f)
}

//小数向下取整保留n位小数
//f 小数
func FloatDownN(f float64, n int) float64 {
	var o float64 = 1
	for i := 0; i < n; i++ {
		o = o * 10
	}
	return math.Floor(f*o) / o
}

func entryptUserIdCard(i string) (ei string) {
	if len(i) != 18 {
		ei = i
	} else {
		midi := i[4:14]
		ei = "****" + midi + "****"
	}
	return ei
}

func entryptUserMobile(p string) (ep string) {
	if len(p) != 11 {
		ep = p
	} else {
		fp := p[:3]
		bp := p[7:]
		ep = fp + "****" + bp
	}
	return ep
}

func EntryptBankCard(p string) (ep string) {
	if len(p) < 6 {
		ep = p
	} else {
		bp := p[len(p)-6:]
		strx := ""
		for i := 0; i < len(p)-6; i++ {
			strx = strx + "*"
		}
		ep = strx + bp
	}
	return ep
}
