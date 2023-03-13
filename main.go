package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "root:password@tcp(127.0.0.1:3306)/ribbit?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqldb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(100)
	sqldb.SetConnMaxLifetime(time.Second * 30)
	AutoMigrate(db)
	CreateUser(db)

}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
}

func CreateUser(db *gorm.DB) {
	users := []User{
		{
			ID:               "USER.ID.1",
			Name:             "USER.NAME.1",
			CreatedAt:        time.Now(),
			Password:         "USER.PASSWORD.1",
			SensorPrivateKey: "PRIVATE.KEY.1",
		},
		{
			ID:               "USER.ID.2",
			Name:             "USER.NAME.2",
			CreatedAt:        time.Now(),
			Password:         "USER.PASSWORD.2",
			SensorPrivateKey: "PRIVATE.KEY.2",
		},
		{
			ID:               "USER.ID.3",
			Name:             "USER.NAME.3",
			CreatedAt:        time.Now(),
			Password:         "USER.PASSWORD.3",
			SensorPrivateKey: "PRIVATE.KEY.3",
		},
	}
	res := db.Create(&users)
	if res.Error != nil {
		panic(res.Error)
	}
	for _, user := range users {
		fmt.Printf("school.ID: %d\n", user.ID)
	}
}
