package db

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type UserDB struct {
	ID            string `gorm:"column:id;primary_key"`
	Name          string `gorm:"column:name"`
	Email         string `gorm:"column:email"`
	CreatedAt     time.Time
	Password      string `gorm:"column:password"`
	ProjectID     string `gorm:"column:project_id"`
	DeviceID      string `gorm:"column:device_id"`
	GoliothAPIKey string `gorm:"column:golioth_api_key"`
	PSK           string `gorm:"column:psk"`
}

type UserRepo struct {
	db *gorm.DB
}
type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	ProjectID string
	DeviceID  string
}

// TableName sets the table name for the UserDB model
func (UserDB) TableName() string {
	return "users"
}

// NewUserRepo initializes a new instance of the [UserRepo] type
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
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

	db.AutoMigrate(&User{})
}

// CreateUser will add a single new user to database
func CreateUser(user User) error {
	user.ID = uuid.New().String()
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetAllUsers retrieves a list of all users and their info from the database
func GetAllUsers() ([]User, error) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// DeleteUserByEmail deletes a user from the database identified by email
func DeleteUserByEmail(email string) error {
	return db.Delete(&User{}, "email = ?", email).Error
}
