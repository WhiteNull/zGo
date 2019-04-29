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
	////sit
	//ip := "139.196.72.249"
	//port := 9090
	//username := "dbtest"
	//password := "Test@Woda2018"
	//dbname := "zxx_sit"
	//local
	ip := "192.168.199.73"
	port := 3306
	username := "root"
	password := "password"
	dbname := "zxx"
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
	item := make([]*PreLeakOutStruct, 0)
	db := gGorm.New().Table(Dao.BillWeeklyBatchPeriodSplitDetail{}.TableName()).
		Select("bill_weekly_batch_period_split_detail.bill_weekly_batch_period_split_detail_id as split_id,bill_weekly_batch_period_split_detail.ent_id,bill_weekly_batch_period_split_detail.trgt_sp_id,"+
			"bill_weekly_batch_detail.srce_sp_id,user_unique.real_name,bill_weekly_batch_detail.work_card_no,"+
			"user.mobile,bill_weekly_batch_period_split_detail.id_card_num,SUM(bill_weekly_batch_period_split_detail.advance_pay_amt) as paid_amt,"+
			"bill_weekly_batch_detail.entry_dt,bill_weekly_batch_detail.leave_dt,user_idcard_audit.audit_tm").
		Joins("LEFT JOIN bill_weekly_batch_detail on bill_weekly_batch_period_split_detail.bill_weekly_batch_detail_id = bill_weekly_batch_detail.bill_weekly_batch_detail_id").
		Joins("LEFT JOIN user_unique on bill_weekly_batch_period_split_detail.uuid = user_unique.uuid").
		Joins("LEFT JOIN user_idcard_audit on bill_weekly_batch_period_split_detail.id_card_num = user_idcard_audit.id_card_num").
		Joins("LEFT JOIN user on user_idcard_audit.user_id = user.user_id").
		Where("bill_weekly_batch_period_split_detail.bill_related_mo =?", "2019-03-01").
		Where("bill_weekly_batch_period_split_detail.srce_sp_audit_sts =?", 2).
		Where("bill_weekly_batch_period_split_detail.trgt_sp_audit_sts =?", 2).
		Where("bill_weekly_batch_period_split_detail.is_deleted =?", 1).
		Where("bill_weekly_batch_detail.entry_dt <= bill_weekly_batch_period_split_detail.end_dt").
		Where("(bill_weekly_batch_detail.leave_dt = '0000-00-00') or (bill_weekly_batch_detail.leave_dt !='0000-00-00' and bill_weekly_batch_detail.leave_dt >= bill_weekly_batch_period_split_detail.begin_dt)").
		Where("user_idcard_audit.audit_sts =?", 2).
		Where("user_idcard_audit.is_deleted =?", 1).
		Where("user.is_deleted = ?", 1).
		Where("user.is_certed = ?", 1)
	db = db.Group("ent_id,trgt_sp_id,srce_sp_id,real_name,work_card_no,mobile,id_card_num,entry_dt,leave_dt")
	db.Find(&item)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			fmt.Println("1111111111")
		}
		fmt.Println("22222222")
	}
	preLeakOutStructMap := map[string]*PreLeakOutStruct{}
	for _, v := range item {
		preLeakOutStruct := preLeakOutStructMap[v.IdCardNum]
		if preLeakOutStruct == nil {
			preLeakOutStructMap[v.IdCardNum] = v
		} else {
			tnew, _ := time.Parse("2006-01-02 15:04:05", v.AuditTm)
			told, _ := time.Parse("2006-01-02 15:04:05", preLeakOutStruct.AuditTm)
			if tnew.After(told) || tnew.Equal(told) {
				if v.SrceSpId > preLeakOutStruct.SplitId {
					preLeakOutStructMap[v.IdCardNum] = v
				}
			}
		}
	}
	newitem := make([]*PreLeakOutStruct, 0)
	for _, v := range preLeakOutStructMap {
		newitem = append(newitem, v)
	}
	fmt.Println("33333333")
}

type PreLeakOutStruct struct {
	EntId      int64
	TrgtSpId   int64
	SrceSpId   int64
	RealName   string
	WorkCardNo string
	Mobile     string
	IdCardNum  string
	PaidAmt    int64
	EntryDt    string
	LeaveDt    string
	AuditTm    string
	SplitId    int64
}
