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

//将日志打印的sql转换成可执行的sql
func main() {
	fmt.Println("sql format start")

	strsql := "UPDATE `user_clock_rec` SET `advance_pay_amt` = ?, `calculate_pay_sts` = ?, `clock_in_is_valid` = ?, `clock_in_sts` = ?, `clock_in_tm` = ?, `clock_in_typ` = ?, `clock_out_is_valid` = ?, `clock_out_sts` = ?, `clock_out_tm` = ?, `clock_out_typ` = ?, `clock_sts` = ?, `name_list_id` = ?, `remark` = ?, `updated_tm` = ?, `uuid` = ?  WHERE (uuid = ? OR user_id = ?) AND (clock_dt = ?) AND (is_deleted = ?) [0 0 1 1 2019-04-13 19:51:33 2 1 1 2019-04-14 07:52:39 2 1 127785  2019-04-15 14:46:10 151746 151746 154921 2019-04-13 1]"

	esql, err := replaceValue(strsql)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("可执行sql:")
		fmt.Println(esql)
	}
	fmt.Println("sql format end")
}
