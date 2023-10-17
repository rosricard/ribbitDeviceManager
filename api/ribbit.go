package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rosricard/userAccess/db"
)

//TODO: store config in db a json string, create config file to handle creation and parsing of these details

const (
	projectID = "ribbit-test-569244"
	baseURL   = "https://api.golioth.io"
	apiKey    = "R7aJE5qW4DNHJTgy9JpbmZYrFXnRTY8S"
)

type Device struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
	PreSharedKey string    `json:"preSharedKey"`
	ProjectID    string    `json:"projectId"`
}

type DeviceList struct {
	List []Device `json:"list"`
}

func CreateUser(c *gin.Context) {
	var userInput struct {
		ID        string
		Name      string
		Email     string
		Password  string
		ProjectID string
		DeviceID  string
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := db.User{
		ID:       userInput.ID,
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: userInput.Password, //TODO: hash password
	}

	if err := db.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func GetAllUsers(c *gin.Context) {
	users, err := db.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func DeleteUser(c *gin.Context) {
	email := c.Param("email")

	if err := db.DeleteUserByEmail(email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// TODO: track user login

// joinUserDevice adds a device to a user
func joinUserDevice(user db.User, device db.Device) error {

	//create device
	d, err := createNewDevice()
	if err != nil {
		return err
	}

	//create private key
	psk, err := createPSK(d.DeviceId)
	if err != nil {
		return err
	}

	dev := Device{
		ID:           d.DeviceId,
		Name:         d.Name,
		CreatedAt:    psk.CreatedAt,
		PreSharedKey: psk.PreSharedKey,
		ProjectID:    d.ProjectID,
	}

	//save device to db

	log.Printf("device: %v", dev)
	//getUser info

	//combine user info and device info

	// add device to db if success was confirmed
	return nil
}

// createDevice creates a new device and returns the device id and psk
func createDevice(c *gin.Context) {
	// create device
	device, err := createNewDevice()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// create private key for device
	pskData, err := createPSK(device.DeviceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"device": device.DeviceId, "psk": pskData.PreSharedKey})

}

// TODO: setup config files with projectID, tagIds, APIkey, etc
// user logs in
// add device to table
// TODO: on app startup, run a check against the golioth API to get all devices and compare against the database

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/createusers", CreateUser)
	//TODO: change this to GetUser
	r.GET("/getusers", GetAllUsers)
	r.DELETE("/users/:email", DeleteUser)
	r.POST("/createDevice", createDevice)
	// get api keys
	//https://api.golioth.io/v1/projects/ribbit-test-569244/apikeys
	return r
}
