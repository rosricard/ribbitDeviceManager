package db

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/rosricard/userAccess/model"
	"gorm.io/gorm"
)

type UserDB struct {
	ID               string `gorm:"column:id;primary_key"`
	Name             string `gorm:"column:name"`
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

// ToModel maps the database user model to the graphQL user model
func (u *UserDB) ToModel() *model.User {
	string_uuid, err := uuid.FromString(u.ID)
	if err != nil {
		fmt.Println(err)
	}

	return &model.User{
		ID:       string_uuid,
		Name:     u.Name,
		Password: u.Password,
	}
}

// getAllUsers requests a list of users from the database
func GetAllUsers(db *gorm.DB) ([]*UserDB, error) {
	var users []*UserDB
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
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
