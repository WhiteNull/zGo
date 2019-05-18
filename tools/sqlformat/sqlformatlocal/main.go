package main

import (
	"fmt"
	"strings"
	"time"
)

//将日志打印的sql转换成可执行的sql
func main() {
	fmt.Println("sql format start")

	strsql := "SELECT name_list.name_list_id as name_ist_id,name_list.created_tm as create_tm,name_list.updated_tm as update_tm,name_list.ent_id as real_ent_id,name_list.entry_dt as entry_dt,name_list.id_card_num as id_card_num,name_list.intv_dt as intv_dt,name_list.intv_sts as intv_sts,name_list.is_valid as is_valid,name_list.leave_dt as leave_dt,user_unique.real_name as realname,name_list.srce_sp_id as srce_sp_id,name_list.trgt_sp_id as trgt_sp_id,name_list.uuid as uuid,name_list.work_sts as work_sts,name_list.work_card_no as work_card_no,user_unique.gender as gender,name_list.rcrt_order_id as rcrt_order_id,name_list.jff_sp_ent_id as work_id,name_list.jff_sp_ent_name as work_name FROM `name_list` left join user_unique on name_list.uuid = user_unique.uuid and name_list.uuid != 0 WHERE (name_list.entry_dt >= ?) AND (name_list.entry_dt < ?) AND (name_list.entry_dt >= ?) AND (name_list.entry_dt < ?) AND (name_list.intv_dt >= ? AND name_list.intv_dt <= ?) AND (name_list.is_deleted = 1) ORDER BY name_list.intv_dt desc,name_list.name_list_id descLIMIT 10 OFFSET 0 [2018-02-01 2018-05-02 2017-10-01 2019-05-04 2018-05-01 00:00:00 2019-05-10 23:59:59]"

	esql, err := replaceValue(strsql)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("可执行sql:")
		fmt.Println(esql)
	}
	fmt.Println("sql format end")
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
