package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rosricard/userAccess/db"
)

const (
	projectID = "ribbit-test-569244"
)

type Device struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Identity     string `json:"identity"`
	CreatedAt    string `json:"createdAt"`
	PreSharedKey string `json:"preSharedKey"`
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

const (
	baseURL = "https://api.golioth.io"
	apiKey  = "R7aJE5qW4DNHJTgy9JpbmZYrFXnRTY8S"
	apiURL  = "https://api.golioth.io/v1/projects/ribbit-test-569244/devices/64194746a946a2ad67aba7ad/credentials"
)

// goliothGetRequest handles GET requests to external APIs
func goliothGetRequest(c *gin.Context) {

	// Create a new HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Add headers to the HTTP request
	req.Header.Set("X-API-Key", apiKey)

	// Make the API call
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var devices DeviceList

	err1 := json.Unmarshal(body, &devices)
	if err1 != nil {
		fmt.Println("Error:", err)
		return
	}

	// return the response to the client
	c.JSON(http.StatusOK, gin.H{"message": resp.Status, "data": devices})
}

//TODO: setup config files with projectID
// user logs in
// create a new device
// add device to table

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/createusers", CreateUser)
	//TODO: change this to GetUser
	r.GET("/getusers", GetAllUsers)
	r.DELETE("/users/:email", DeleteUser)
	r.GET("/goliothGetRequest", goliothGetRequest)
	//TODO: Add single device API
	//can create a new "blank" device
	//r.GET("/v1/projects/ribbit-test-569244/devices", getAllDevices)
	//get devices from golioth
	//r.GET("/devices", getAllDevices)
	//get the user device identity and PSK
	//r.GET("/v1/projects/ribbit-test-569244/credentials", getAllDevices)
	// get api keys
	//https://api.golioth.io/v1/projects/ribbit-test-569244/apikeys

	//TODO: generate credentials for the user

	return r
}
