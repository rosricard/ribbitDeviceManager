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
