// package db provides functionality for interacting with a relational database
package db

import (
	"errors"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type Repository struct {
	users   *UserRepo
	devices *DeviceRepo
	// TODO: add in zap logger
}

// AutoMigrate will automatically migrate the database and correct schema errors on startup
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Device{})
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &Repository{
		users:   NewUserRepo(db),
		devices: NewDeviceRepo(db),
	}, nil
}

// ConnectDatabase initalizes the sql database connection and gorm
func ConnectDatabase() {
	//TODO: setup db connection as an env variable
	dsn := "root:billybob123@tcp(127.0.0.1:3306)/ribbit?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	db.AutoMigrate(&User{}, &Device{})
}

func (r *Repository) Users() *UserRepo {
	return r.users
}

func (r *Repository) Devices() *DeviceRepo {
	return r.devices
}
