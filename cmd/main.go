package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Database func() *gorm.DB
	database *gorm.DB

	Logger func() *log.Logger
	logger *log.Logger
}

type Repository struct {
	users *UserRepo
}

type user struct {
	ID               string `gorm:"column:id;primary_key"`
	name             string `gorm:"column:name"`
	password         string `gorm:"column:password"`
	sensorPrivateKey string `gorm:"column:private_key"`
}

type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo initializes a new instance of the [UserRepo] type
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

// Automigrate automatically updates any configured database table if a mis-match in config is detected
func Automigrate(db *gorm.DB) error {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:password@tcp(127.0.0.1:3306)/ribbit?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&user{})
	return err
}

// NewRepository initializes a new instance of the [Repository] type
//
// [Automigrate] will be called to update the db as part of initialization
func NewRepository(db *gorm.DB) (*Repository, error) {
	if err := Automigrate(db); err != nil {
		return nil, err
	}
	return &Repository{
		users: NewUserRepo(db),
	}, nil
}

func main() {
	config := new(Config)
	_, err := NewRepository(config.Database())
	if err != nil {
		//Config.logger.Fatal(err)
		print(err)
	}
}
