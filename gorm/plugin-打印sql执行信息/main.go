package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"test/gorm/src"
)


type userInfo struct {
	Id int64 `gorm:"column:id"`
	username string `gorm:"column:username"`
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:123456@tcp(127.0.0.1:3306)/test2?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil{
		log.Fatal(err)
	}


	//使用插件
	err = db.Use(&src.TracePlugin{})
	if err != nil{
		log.Fatal(err)
	}

	user := []userInfo{}

	err = db.Table("user").Model(&userInfo{}).Where("id < ?",10).Find(&user).Error

	if err != nil{
		log.Fatal(err)
	}
}


