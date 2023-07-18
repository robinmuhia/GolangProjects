package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var(
	DBConn *gorm.DB
)

func Connect(){
	d, err := gorm.Open("mysql","robinmuhia:qwertyuiop@tcp(localhost)/simplerest?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		panic(err)
	}
	DBConn = d
}

func GetDB() *gorm.DB{
	return DBConn
}
