package db

import (
	"log"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type UserDB struct {
	ID               string `gorm:"column:id;primary_key"`
	Name             string `gorm:"column:name"`
	Email            string `gorm:"column:email"`
	CreatedAt        time.Time
	Password         string `gorm:"column:password"`
	SensorPrivateKey string `gorm:"column:private_key"`
}

type UserRepo struct {
	db *gorm.DB
}

// TableName sets the table name for the UserDB model
func (UserDB) TableName() string {
	return "users"
}

// NewUserRepo initializes a new instance of the [UserRepo] type
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

type User struct {
	ID         uuid.UUID
	Name       string
	Email      string
	Password   string
	PrivateKey string
}

func ConnectDatabase() {
	dsn := "root:password@tcp(127.0.0.1:3306)/ribbit?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	db.AutoMigrate(&User{})
}

func CreateUser(id uuid.UUID, name, email, password, pk string) error {
	user := User{
		ID:         id,
		Name:       name,
		Email:      email,
		Password:   password, // TODO: hash password
		PrivateKey: pk,       // TODO: retrieve private key from golioth api
	}
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAllUsers() ([]User, error) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func DeleteUser(id uint) error {
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
