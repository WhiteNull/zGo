package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func main() {
	fileurl := "/Users/gongzhihong/Documents/OneDrive/files online/sql.sql"
	fmt.Println("format sql start")
	fmt.Println("file url is:" + fileurl)
	b, err := ioutil.ReadFile(fileurl)
	if err != nil {
		fmt.Println("io real file err:")
		fmt.Print(err)
		return
	}
	strsql := string(b)

	rv, err := replaceValue(strsql)
	if err != nil {
		fmt.Println("replace ? to value err:")
		fmt.Println(err)
	}
	bw := []byte(strsql + "\r\r\r" + rv)
	err = ioutil.WriteFile(fileurl, bw, 0644)
	if err != nil {
		fmt.Print("io write file err:")
		fmt.Print(err)
		return
	} else {
		fmt.Println("success")
	}
	fmt.Println("format sql end")
}

func replaceValue(strsql string) (string, error) {
	s := strings.Index(strsql, "[")
	e := strings.LastIndex(strsql, "]")
	if s <= 0 || e <= 0 {
		return "", fmt.Errorf("未找到[]中的值")
	}
	sql := strsql[:s]
	sqlvalue := strsql[s+1 : e]
	values := strings.Split(sqlvalue, " ")
	//判断sql语句中?的数量
	countq := strings.Count(sql, "?")
	if countq == len(values) {
		for i := 0; i < len(values); i++ {
			value := values[i]
			sql = strings.Replace(sql, "?", value, 1)
		}
	} else if len(values) > countq {
		//一个时间中可能含有空格导致分割值出错
		for i := 0; i < len(values); i++ {
			value := values[i]
			if i < len(values)-1 {
				nextvalue := values[i+1]
				timev := value + " " + nextvalue
				_, err := time.Parse("2006-01-02 15:04:05", timev)
				if err != nil {
					sql = strings.Replace(sql, "?", value, 1)
				} else {
					sql = strings.Replace(sql, "?", timev, 1)
					i++
				}
			} else {
				sql = strings.Replace(sql, "?", value, 1)
			}
		}
	} else {
		return "", fmt.Errorf("[]中的值数量比sql中的?数量少")
	}
	return sql, nil
}
