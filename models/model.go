package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	//"github.com/go-sql-driver/mysql"
	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	//_ "github.com/jinzhu/gorm/dialects/postgres" //pgSQL
)

var GlobalDb *gorm.DB

// User 用户表
type User struct {
	//gorm.Model
	ID       uint   `gorm:"primary_key"`
	Phone    string `gorm:"size:100;unique;not null"`
	Password string `gorm:"size:100;not null"`
	Avatar   string `gorm:"size:1024;default:'https://s2.ax1x.com/2019/05/25/Vknme1.png'"`
	//NickName    string `gorm:"size:255;default:'默认用户'"`
	RealName string `gorm:"size:255;default:'默认用户'"`
	Role     string `gorm:"size:16;default:'user'"` //user editor admin
	Lock     string `gorm:"size:16;default:'unlock'"`
	Email    string `gorm:"size:100"`

	Car   Car //关联车号，以车号来区分  tips:一定要不保存关联，否则原始数据会被覆盖
	CarID uint

	LastLoginAt uint
	LastLoginIP string `gorm:"size:32"`
}

// Place 位置表
type Place struct {
	//gorm.Model
	ID       uint   `gorm:"primary_key"`
	Position string `gorm:"size:255;unique;not null"`
	Remarks  string `gorm:"size:255"`
}

// Car 车号 DM05等
type Car struct {
	//gorm.Model
	ID      uint   `gorm:"primary_key"`
	Car     string `gorm:"size:255;unique;not null"`
	Remarks string `gorm:"size:255"`
}

// Material 材料表
type Material struct {
	//gorm.Model
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"size:255;not null"` //名称
	Model    string `gorm:"size:255"`          //型号
	NickName string `gorm:"size:255"`          //俗称
	Unit     string `gorm:"size:255"`          //计量单位

	//Car   Car //关联车号，以车号来区分  tips:一定要不保存关联，否则原始数据会被覆盖
	//CarID uint

	Place    Place //关联货架  tips:一定要不保存关联，否则原始数据会被覆盖
	PlaceID  uint
	Floor    int //层
	Location int //位

	Count        int    //数量
	PrepareCount int    //常备数量
	WarnCount    int    //警报数量
	Marks        string `gorm:"size:255"` //备注

	User     User //创建用户
	UserID   uint
	CreateAt int64 //创建时间
}

// Record  记录表
type Record struct {
	//gorm.Model
	ID         uint `gorm:"primary_key"`
	Material   Material
	MaterialID uint

	User   User //创建用户
	UserID uint

	//Car   Car //关联车号，以车号来区分  tips:一定要不保存关联，否则原始数据会被覆盖
	//CarID uint

	CreateAt    int64 //创建时间
	UpdateAt    int64
	Type        string //记录类型 "receive":领料  "send":发料
	CountChange int    //变动数量 如1  5
	BeforeCount int    //变动前数量
	AfterCount  int    //变动后数量
	Marks       string `gorm:"size:255"` //备注
}

//初始化
func init() {
	// 连接数据库配置
	dsn := "cailiao:123456@tcp(127.0.0.1:3306)/cailiao?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//连接池
	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	//db, err := gorm.Open("mysql", "cailiao:123456@tcp(127.0.0.1:3306)/cailiao?charset=utf8mb4&parseTime=True&loc=Local")
	//defer db.Close()

	GlobalDb = db

	// 表迁移
	//err = GLOBAL_DB.AutoMigrate(&User{}, &Place{}, &Dog{})
	err = GlobalDb.AutoMigrate(&User{}, &Place{}, &Car{}, &Material{}, &Record{})
	if err != nil {
		fmt.Println("表迁移失败", err)
		panic(err)
	}

	fmt.Println("数据库初始化成功")
}
