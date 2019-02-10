package main

import (
	"fmt"
	"strings"
)

func replaceValue(strsql string) (string, error) {
	s := strings.Index(strsql, "[")
	e := strings.LastIndex(strsql, "]")

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

	strsql := "SELECT remittance_app_id FROM `remittance_app`  WHERE (  bill_batch_id = ? AND bill_batch_typ = ? AND auth_sts !=? AND is_deleted = ?) LIMIT 1 [12580 2 3 1] 1"

	esql, err := replaceValue(strsql)
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println("可执行sql:")
		fmt.Println(esql)
	}
	fmt.Println("sql format end")
}
