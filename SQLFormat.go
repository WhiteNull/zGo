package main

import (
	"fmt"
	"strings"
)

func replaceValue(strsql string) (string, error) {
	s := strings.Index(strsql, "[")
	e := strings.LastIndex(strsql, "]")
	if s <= 0 || e <= 0 {
		return "", fmt.Errorf("要格式化的SQL不正确")
	}
	sql := strsql[:s]
	sqlvalue := strsql[s+1 : e]
	values := strings.Split(sqlvalue, " ")
	for i := 0; i < len(values); i++ {
		value := values[i]
		sql = strings.Replace(sql, "?", value, 1)
	}
	return sql, nil
}

func main() {
	fmt.Println("sql format start")

	strsql := ""

	esql, err := replaceValue(strsql)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("可执行sql:")
		fmt.Println(esql)
	}
	fmt.Println("sql format end")
}
