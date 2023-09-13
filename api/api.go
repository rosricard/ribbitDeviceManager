package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rosricard/userAccess/db"
)

const (
	projectID = "ribbit-test-569244"
)

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

	// fmt.Println(resp.Body)
	c.JSON(http.StatusOK, gin.H{"message": string(body)})
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
