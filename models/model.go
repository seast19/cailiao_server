package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //pgSQL
)

// User 用户表
type User struct {
	gorm.Model
	Phone       string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	Avatar      string `gorm:"default:'https://s2.ax1x.com/2019/05/25/Vknme1.png'"`
	Nickname    string `gorm:"default:'默认用户'"`
	Role        string `gorm:"default:'user'"`
	Lock        string `gorm:"default:'unlock'"`
	Email       string
	LastLoginAt int64
	LastLoginIP string
}

func init() {
	// 连接数据库
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=cailiao dbname=cailiao password=123456 sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 表迁移
	db.AutoMigrate(&User{})

	fmt.Println("数据库初始化成功")
}

func getConn() (*gorm.DB, error) {
	// 连接数据库
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=cailiao dbname=cailiao password=123456 sslmode=disable")
	if err != nil {
		fmt.Printf("连接数据库失败 -> %s\n", err)
		return nil, err
	}

	fmt.Println("连接数据库成功")
	return db, nil
}
