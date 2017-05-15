package dao

import (
	"../myconstant"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func CheckTableAndColumn(userName string, pwd string) {
	db, err := gorm.Open(myconstant.MYSQL, userName+":"+pwd+"@/douban?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	if (db.HasTable(&Movie{})) {
		db.AutoMigrate(&Movie{}) //当你在module 中添加了新的字段后,会同步到表中,
	} else {
		db.Set("gorm:table_options", "ENGINE=InnoDB,charset=utf8").
			CreateTable(&Movie{})
	}
}

func WriteMovie(movies ...Movie) {
	db, err := gorm.Open(myconstant.MYSQL, "root:3252860@/douban?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		panic(err)
	}
	//db.Exec()
	for _, moive := range movies {
		db := db.Create(moive)
		fmt.Println(db.Value)
	}

}
