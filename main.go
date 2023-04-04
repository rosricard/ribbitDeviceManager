package main

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "root:password@tcp(127.0.0.1:3306)/ribbit?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	// create web server
	srv := NewServer()
	http.ListenAndServe(":8080", srv)

	// gorm for db read / write
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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("hello, world")
}
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&UserDB{})
}

func CreateUser(db *gorm.DB) {
	users := []UserDB{
		{
			ID:               "USER_ID.1",
			Name:             "USER_NAME.1",
			CreatedAt:        time.Now(),
			Password:         "USER_PASSWORD.1",
			SensorPrivateKey: "PRIVATE_KEY.1",
		},
		{
			ID:               "USER_ID.2",
			Name:             "USER_NAME.2",
			CreatedAt:        time.Now(),
			Password:         "USER_PASSWORD.2",
			SensorPrivateKey: "PRIVATE_KEY.2",
		},
		{
			ID:               "USER_ID.3",
			Name:             "USER_NAME.3",
			CreatedAt:        time.Now(),
			Password:         "USER_PASSWORD.3",
			SensorPrivateKey: "PRIVATE_KEY.3",
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
