package main

import (
	"fmt"
	"git.woda.ink/zxx/common/utils/Dao"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
)

var gGorm *gorm.DB

var dbPingOnce sync.Once

func BuildConnInfo(user string, pwd string, ip string, port string, db string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, pwd, ip, port, db)
}

func dbTimerPing(ggorm *gorm.DB) {
	tick := time.NewTicker(time.Second * 300)

	for {
		select {
		case <-tick.C:
			ggorm.DB().Ping()
		}
	}
}

func main() {
	//连接db
	////alpha
	//ip := "139.224.170.14"
	//port := 10133
	//username := "testing"
	//password := "Test!@#321"
	//dbname := "zxx"
	//sit
	ip := "139.196.72.249"
	port := 9090
	username := "dbtest"
	password := "Test@Woda2018"
	dbname := "zxx_sit"
	gGorm, err := gorm.Open("mysql", BuildConnInfo(username, password, ip, fmt.Sprintf("%d", port), dbname))
	if err != nil {
		fmt.Println("InitDb", "gorm.Open", err)
	}
	err = gGorm.DB().Ping()
	if err != nil {
		fmt.Println("gORM.DB().Ping", "gorm.Open", err)
	}
	gGorm.DB().SetMaxIdleConns(5)
	gGorm.DB().SetMaxOpenConns(10)
	gGorm.DB().SetConnMaxLifetime(time.Hour)
	dbPingOnce.Do(func() {
		go dbTimerPing(gGorm)
	})
	gGorm.LogMode(true)
	//连接db end

	//start do
	//param
	req := _QueryRecordParams{}
	req.RecordIndex = 0
	req.RecordSize = 10
	req.Mobile = "15713774718"
	var count int64
	exportflag := false
	list := make([]ClockRecord, 0)
	//sql
	db := gGorm.New().Table(Dao.UserClockRec{}.TableName()).
		Select("user_unique.real_name," +
			"user_unique.id_card_num," +
			"user.mobile," +
			"name_list.work_card_no," +
			"name_list.intv_dt as interview_dt," +
			"ent.ent_short_name as ent_name," +
			"tsp.sp_short_name as trgt_sp_name," +
			"ssp.sp_short_name as srce_sp_name," +
			"user_clock_rec.clock_dt," +
			"user_clock_rec.clock_in_tm," +
			"user_clock_rec.clock_in_addr," +
			"user_clock_rec.clock_out_tm," +
			"user_clock_rec.clock_out_addr," +
			"user_clock_rec.advance_pay_amt as amount," +
			"user_clock_rec.clock_out_typ as is_clock_out_repaired," +
			"name_list.srce_sp_id as srce_sp_id," +
			"user_clock_rec.clock_out_sts," +
			"ent.ent_id as ent_id," +
			"user_clock_rec.remark," +
			"user_clock_rec.clock_in_sts," +
			"user_clock_rec.user_clock_rec_id as clock_rec_id," +
			"name_list.trgt_sp_id as trgt_sp_id," +
			"user_clock_rec.clock_in_typ as is_clock_in_repaired," +
			"user_clock_rec.calculate_pay_sts").
		Joins("LEFT JOIN user on user_clock_rec.user_id = user.user_id " +
			"and user.is_deleted = 1").
		Joins("LEFT JOIN user_unique on user.uuid = user_unique.uuid " +
			"and user_unique.is_deleted = 1").
		Joins("LEFT JOIN name_list on user_clock_rec.name_list_id = name_list.name_list_id " +
			"and name_list.is_deleted = 1").
		Joins("LEFT JOIN ent on name_list.ent_id = ent.ent_id " +
			"and ent.is_deleted = 1").
		Joins("LEFT JOIN sp tsp on name_list.trgt_sp_id = tsp.sp_id").
		Joins("LEFT JOIN sp ssp on name_list.srce_sp_id = ssp.sp_id").
		Order("user_clock_rec.user_clock_rec_id DESC")
	if len(req.ClockStartDt) > 0 {
		db = db.Where("user_clock_rec.clock_dt >= ?", req.ClockStartDt)
	}
	if len(req.ClockEndDt) > 0 {
		db = db.Where("user_clock_rec.clock_dt <= ?", req.ClockEndDt)
	}
	if req.EntID != 0 && req.EntID != -9999 {
		db = db.Where("ent.ent_id = ?", req.EntID)
	}
	if len(req.IDCardNum) > 0 {
		db = db.Where("user_unique.id_card_num = ?", req.IDCardNum)
	}
	if req.IsRepaired != 0 && req.IsRepaired != -9999 {
		if req.IsRepaired == 1 {
			db = db.Where("user_clock_rec.clock_in_typ in (2,3) and user_clock_rec.clock_out_typ in (2,3)")
		}
		if req.IsRepaired == 2 {
			db = db.Where("user_clock_rec.clock_in_typ = 1 and user_clock_rec.clock_out_typ = 1")
		}
	}
	if req.IsWeekly != 0 && req.IsWeekly != -9999 {
		if req.IsWeekly == 1 {
			db = db.Where("user_clock_rec.advance_pay_amt > 0")
		}
		if req.IsWeekly == 2 {
			db = db.Where("user_clock_rec.advance_pay_amt = 0")
		}
	}
	if len(req.Mobile) > 0 {
		db = db.Where("user.mobile = ?", req.Mobile)
	}
	if len(req.RealName) > 0 {
		db = db.Where("user_unique.real_name = ?", req.RealName)
	}
	if req.SrceSpID != 0 && req.SrceSpID != -9999 {
		db = db.Where("ssp.sp_id = ?", req.SrceSpID)
	}
	if req.TrgtSpID != 0 && req.TrgtSpID != -9999 {
		db = db.Where("tsp.sp_id = ?", req.TrgtSpID)
	}
	if len(req.WorkCardNo) > 0 {
		db = db.Where("name_list.work_card_no = ?", req.WorkCardNo)
	}
	db = db.Where("user_clock_rec.is_deleted = ?", 1)
	countdb := db
	if exportflag == true {
		db = db.Find(&list)
	} else {
		db = db.Offset(req.RecordIndex).Limit(req.RecordSize).Find(&list)
	}
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			fmt.Println("111111111111")
		}
		fmt.Println("2222222")
		fmt.Println(db.Error.Error())
	}
	//get count
	countitem := CountStruct{}
	countdb = countdb.Select("count(*) as cs").Find(&countitem)
	count = countitem.Cs
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			fmt.Println("111111111111")
		}
		fmt.Println("2222222")
		fmt.Println(countdb.Error.Error())
	}
	for i := 0; i < len(list); i++ {
		list[i].IsRepaired = "否"
		if list[i].IsClockInRepaired == 1 {
			list[i].IsClockInRepaired = 2
		} else {
			list[i].IsClockInRepaired = 1
			list[i].IsRepaired = "是"
		}
		if list[i].IsClockOutRepaired == 1 {
			list[i].IsClockOutRepaired = 2
		} else {
			list[i].IsClockOutRepaired = 1
			list[i].IsRepaired = "是"
		}
	}
	fmt.Println("count")
	fmt.Println(count)
	fmt.Println("list")
	fmt.Println(list)
}

type _QueryRecordParams struct {
	RealName     string //姓名
	WorkCardNo   string //工号
	ClockStartDt string //开始打卡日期
	IDCardNum    string //身份证号
	TrgtSpID     int    //劳务ID
	EntID        int    //企业ID
	IsWeekly     int    //是否周薪
	Mobile       string //手机号
	SrceSpID     int    //中介ID
	ClockEndDt   string //结束打卡日期
	IsRepaired   int    //是否补卡
	RecordIndex  int64  //取记录起始数量
	RecordSize   int64  //获取记录数量
}

type ClockRecord struct {
	RealName           string `TExcel:"姓名"`   ///姓名
	IDCardNum          string `TExcel:"身份证号"` //身份证
	Mobile             string `TExcel:"手机号码"` //手机号
	WorkCardNo         string `TExcel:"工号"`   ///工号
	InterviewDt        string `TExcel:"面试日期"` //面试日期
	EntName            string `TExcel:"企业"`   //企业名称
	TrgtSpName         string `TExcel:"劳务"`   //劳务名称
	SrceSpName         string `TExcel:"中介"`   ///中介名称
	ClockDt            string `TExcel:"打卡日期"` //打卡日期
	ClockInTm          string `TExcel:"上班时间"` //上班打卡时间
	ClockInAddr        string `TExcel:"位置"`   //上班打卡位置
	ClockOutTm         string `TExcel:"下班时间"` //下班打卡时间
	ClockOutAddr       string `TExcel:"位置"`   ///下班打卡位置
	Amount             int64  `TExcel:"金额"`   //金额
	IsRepaired         string `TExcel:"补卡"`   //是否补卡(导出时使用)
	IsClockOutRepaired int64  //是否下班补卡
	SrceSpID           int64  //中介ID
	ClockOutSts        int64  //下班打卡状态
	EntID              int64  //企业ID
	Remark             string //备注
	ClockInSts         int64  //上班打卡状态
	ClockRecID         int64  //打卡记录ID
	TrgtSpID           int64  //劳务ID
	IsClockInRepaired  int64  //是否上班补卡
	CalculatePaySts    string //算钱状态
}

type CountStruct struct {
	Cs int64
}
