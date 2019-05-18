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
	////sit
	//ip := "139.196.72.249"
	//port := 9090
	//username := "dbtest"
	//password := "Test@Woda2018"
	//dbname := "zxx_sit"
	////local
	//ip := "192.168.199.73"
	//port := 3306
	//username := "root"
	//password := "password"
	//dbname := "zxx"
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
	cfgdtlids := []int64{}
	dictid := 0

	item := make([]Dao.CfgDictDetail, 0)
	db := gGorm.New().Table(Dao.CfgDictDetail{}.TableName())
	if len(cfgdtlids) > 0 {
		db = db.Where("cfg_dict_detail_id in (?)", cfgdtlids)
	}
	if dictid != 0 {
		db = db.Where("cfg_dict_id = ?", dictid)
	}
	db = db.Where("is_deleted = ?", 1).
		Find(&item)
	if db.Error != nil {
		fmt.Println("fail")
	} else {
		fmt.Println("success")
		fmt.Println(item)
	}
}

func getNameList(gGorm *gorm.DB, b, e string, entid, trgtid int64) ([]RetentionRateDtl, error) {
	te, _ := time.Parse("2006-01-02", e)
	qe := te.AddDate(0, 0, 1).Format("2006-01-02")
	itemls := make([]RetentionRateDtl, 0)
	db := gGorm.New().Table(Dao.NameList{}.TableName()).
		Select("name_list.entry_dt as date," +
			"name_list.leave_dt," +
			"name_list.ent_id," +
			"ent.ent_short_name," +
			"name_list.trgt_sp_id," +
			"sp.sp_short_name as trgt_sp_short_name," +
			"count(name_list_id) as new_intv_count").
		Joins("LEFT JOIN ent on name_list.ent_id = ent.ent_id " +
			"and ent.is_deleted = 1").
		Joins("LEFT JOIN sp on name_list.trgt_sp_id = sp.sp_id " +
			"and sp.is_deleted = 1").
		Where("name_list.entry_dt >= ?", b).
		Where("name_list.entry_dt < ?", qe).
		Where("name_list.ent_id = ?", entid).
		Where("name_list.trgt_sp_id = ?", trgtid).
		Where("name_list.is_valid = ?", 1).
		Where("name_list.is_deleted = ?", 1).
		Group("name_list.entry_dt").
		Order("name_list.entry_dt").
		Find(&itemls)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, db.Error
	}
	return itemls, nil
}

type RetentionRateDtl struct {
	DataId            int64
	Date              string
	EntId             int64
	EntShortName      string
	NewIntvCount      int64
	RetentionRateList []string
	TrgtSpId          int64
	TrgtSpShortName   string
	LeaveDt           string //前端不用
}

func GetAllSpMap(db *gorm.DB) (map[int64]string, error) {
	item := make([]Dao.Sp, 0)
	db = db.New().Table(Dao.Sp{}.TableName()).
		Select("sp_id,sp_short_name").
		Where("is_deleted = ?", 1).
		Find(&item)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, db.Error
	}
	ret := map[int64]string{}
	for _, v := range item {
		ret[v.SpId] = v.SpShortName
	}
	return ret, nil
}

func GetAllEntMap(db *gorm.DB) (map[int64]string, error) {
	item := make([]Dao.Ent, 0)
	db = db.New().Table(Dao.Ent{}.TableName()).
		Select("ent_id,ent_short_name").
		Where("is_deleted = ?", 1).
		Find(&item)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, db.Error
	}
	ret := map[int64]string{}
	for _, v := range item {
		ret[v.EntId] = v.EntShortName
	}
	return ret, nil
}
