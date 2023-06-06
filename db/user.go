package db

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
func (ur *UserRepo) GetAllUsers(db *gorm.DB) ([]*UserDB, error) {
	var users []*UserDB
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// TODO: handle batching inserts

// Create will add a single user to the database
func (ur *UserRepo) Create(id string, user *model.User) *gorm.DB {
	usr := &UserDB{
		ID:               id,
		Name:             user.Name,
		CreatedAt:        time.Now(),
		Password:         user.Password, //TODO: hash password
		SensorPrivateKey: "testing",     //TODO : retrieve private key from goliath api
	}
	result := ur.db.Create(&usr)
	if result.Error != nil {
		panic(result.Error) // TODO: implement error handling
	}
	return result
}

// DeleteUserByEmail will remove a single user from the database
func (ur *UserRepo) DeleteUserByEmail(email string) *gorm.DB {
	result := ur.db.Delete(&UserDB{}, "email = ?", email)
	if result.Error != nil {
		panic(result.Error) // TODO: implement error handling
	}
	return result
}

// getGoliothPK will retrieve a private key from the Golioth API
func (ur *UserRepo) getGoliothPK(deviceID string) {
	// Golioth API endpoint to request a device private key
	apiURL := "https://api.golioth.io/v1/devices/" + deviceID + "/keys/private"

	// Create a new HTTP client
	client := http.DefaultClient

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers for the request (replace {your_token} with your Golioth API token)
	req.Header.Set("Authorization", "Bearer {your_token}")

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	// Read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status code:", resp.StatusCode)
		fmt.Println("Response:", string(body))
		return
	}

	// Process the response data
	fmt.Println("Private Key:", string(body))
}
