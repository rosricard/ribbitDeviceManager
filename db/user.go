package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

// AutoMigrate will automatically migrate the database and correcrt schema errors on startup
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&UserDB{})
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
