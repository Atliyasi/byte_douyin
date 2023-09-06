package dao

import (
	"douyin_demo/log"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var _db *gorm.DB

func InitDB() {
	USERNAME := "root"
	PASSWORD := "726400Sb."
	DATABASE := "byte_demo"
	PORT := "3306"
	IP := "60.205.176.228"
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb3&parseTime=True&loc=Local",
		USERNAME, PASSWORD, IP, PORT, DATABASE,
	)
	var err error
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.SetLog().Error("[数据库连接]  " + err.Error())
	}
	sqlDB, err := _db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 创建对应数据库
	err = _db.AutoMigrate(&User{}, &VideoUser{}, &Video{}, &FavoriteList{}, &Comment{}, &Relation{}, &Message{})
	//err = _db.AutoMigrate(&User{}, &VideoUser{})
	log.SetLog().Info("[数据库表]  检查成功")
	if err != nil {
		log.SetLog().Error("[数据库表]  " + err.Error())
	}
}

func GetDB() *gorm.DB {
	return _db
}
