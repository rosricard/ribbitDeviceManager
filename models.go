package main

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID               string `gorm:"column:id;primary_key"`
	Name             string `gorm:"column:name"`
	CreatedAt        time.Time
	Password         string `gorm:"column:password"`
	SensorPrivateKey string `gorm:"column:private_key"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	fmt.Printf("User BeforeSave: %+v\n", u)
	return
}
