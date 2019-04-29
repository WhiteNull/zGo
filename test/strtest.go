package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

func main() {
	aa := "123456789012345678"
	fmt.Println("****" + aa[4:14] + "****")

	bb := "1234567890123"
	fmt.Println(bb[0:3] + "****" + bb[len(bb)-4:])

	cc := "2018-08-06 07:14:31"
	fmt.Println(len(cc))
	fmt.Println(cc[0:19])

	dd := make([]interface{}, 0)
	dd = append(dd, "aaaaa")
	dd = append(dd, "bbbbb")
	dd = append(dd, "ccccc")
	dd = append(dd, "ddddd")
	dd = append(dd, "eeeee")
	dd = append(dd, "fffff")

	Write(dd)
}

//一次写一行
func Write(values []interface{}) {
	e := new(Excel)
	row := e.Sheet.AddRow()
	for _, cellValue := range values {
		cell := row.AddCell()
		cell.SetValue(cellValue)
	}
}

type Excel struct {
	File  *xlsx.File
	Sheet *xlsx.Sheet
}
