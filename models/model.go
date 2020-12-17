package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //pgSQL
)

// User 用户表
type User struct {
	ID       uint   `gorm:"primary_key"`
	Phone    string `gorm:"size:100;unique;not null"`
	Password string `gorm:"size:100;not null"`
	Avatar   string `gorm:"size:1024;default:'https://s2.ax1x.com/2019/05/25/Vknme1.png'"`
	//NickName    string `gorm:"size:255;default:'默认用户'"`
	RealName    string `gorm:"size:255;default:'默认用户'"`
	Role        string `gorm:"size:16;default:'user'"` //user editor admin
	Lock        string `gorm:"size:16;default:'unlock'"`
	Email       string `gorm:"size:100"`
	LastLoginAt uint
	LastLoginIP string `gorm:"size:32"`
}

// Place 位置表
type Place struct {
	ID       uint   `gorm:"primary_key"`
	Position string `gorm:"size:255;unique;not null"`
	Remarks  string `gorm:"size:255"`
}

// Material 材料表
type Material struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"size:255;not null"` //名称
	Model    string `gorm:"size:255"`          //型号
	NickName string `gorm:"size:255"`          //俗称
	Unit     string `gorm:"size:255"`          //计量单位

	Place    Place `gorm:"ForeignKey:PlaceID;save_associations:false"` //关联货架  tips:一定要不保存关联，否则原始数据会被覆盖
	PlaceID  uint
	Floor    int //层
	Location int //位

	Count        int    //数量
	PrepareCount int    //常备数量
	WarnCount    int    //警报数量
	Marks        string `gorm:"size:255"` //备注

	User     User `gorm:"ForeignKey:UserID;save_associations:false"` //创建用户
	UserID   uint
	CreateAt int64 //创建时间
}

// record 记录表
type Record struct {
	ID          uint     `gorm:"primary_key"`
	Material    Material `gorm:"ForeignKey:MaterialID;save_associations:false"`
	MaterialID  uint
	User        User `gorm:"ForeignKey:UserID;save_associations:false"` //创建用户
	UserID      uint
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
	// 连接数据库
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=cailiao dbname=cailiao password=123456 sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 表迁移
	db.AutoMigrate(&User{}, &Place{}, &Material{}, &Record{})

	fmt.Println("数据库初始化成功")
}

func getConn() (*gorm.DB, error) {
	// 连接数据库
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=cailiao dbname=cailiao password=123456 sslmode=disable")
	if err != nil {
		fmt.Printf("连接数据库失败 -> %s\n", err)
		return nil, err
	}

	//fmt.Println("连接数据库成功")
	return db, nil
}
