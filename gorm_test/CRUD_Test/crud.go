package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Project struct {
	ID    uint   `gorm:"primarykey"`
	Code  string `gorm:"column: code"`
	Price uint   `gorm:"column: price"`
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:@tcp(127.0.0.1:3306)/gorm?&charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Project{})
	a := &Project{Code: "D42"}
	//Create
	res := db.Create(a) //db.Create属于链式掉用，返回插入的主键ID
	fmt.Println(res.Error)
	fmt.Println(a.ID)
	//	多条插入
	projects := []*Project{{Code: "D40"}, {Code: "D41"}, {Code: "D42"}}
	res = db.Create(projects)
	fmt.Println(res.Error)
	for _, p := range projects {
		fmt.Println(p.ID)
	}
	p := Project{Code: "D42"}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&p)
	//Read

}
