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
	//alpha
	ip := "139.224.170.14"
	port := 10133
	username := "testing"
	password := "Test!@#321"
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
	db := gGorm.New().Table(Dao.BillWeeklyBatchImportDetail{}.TableName()).
		Where("bill_weekly_batch_import_detail_id = ?", 298).
		Update("ent_id", 7)
	if db.Error != nil {
		fmt.Println(",,,,,,,,,")
	} else {
		fmt.Println("lllllllll")
	}

}
