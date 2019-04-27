package main

import (
	"fmt"
	"strings"
	"time"
)

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

//将日志打印的sql转换成可执行的sql
func main() {
	fmt.Println("sql format start")

	strsql := "UPDATE `user_clock_rec` SET `advance_pay_amt` = ?, `calculate_pay_sts` = ?, `clock_out_addr` = ?, `clock_out_device` = ?, `clock_out_ent_loc_id` = ?, `clock_out_is_valid` = ?, `clock_out_latitude` = ?, `clock_out_longitude` = ?, `clock_out_sts` = ?, `clock_out_tm` = ?, `clock_out_typ` = ?, `clock_pay_amt` = ?, `clock_sts` = ?, `credit_subsidy_amt` = ?, `name_list_id` = ?, `uuid` = ?  WHERE `user_clock_rec`.`user_clock_rec_id` = ? [2231 5 222222222 设备Id 0 1 19211 19211 2 2019-04-27 10:51:21 1 2232 2 2233 16223 251722 100000120]"

	esql, err := replaceValue(strsql)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("可执行sql:")
		fmt.Println(esql)
	}
	fmt.Println("sql format end")
}
