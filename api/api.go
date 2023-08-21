package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rosricard/userAccess/db"
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
	apiKey  = "some-api-key"
)

// HandleGetRequest handles GET requests to external APIs
func handleGetRequest(c *gin.Context) {
	url := fmt.Sprintf("%s%s", baseURL, c.Request.RequestURI)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create request")
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to execute request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.String(resp.StatusCode, "Failed to fetch data")
		return
	}

	// Handle the response body here

	c.String(http.StatusOK, "GET request successful")
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/createusers", CreateUser)
	//TODO: change this to GetUser
	r.GET("/getusers", GetAllUsers)
	r.DELETE("/users/:email", DeleteUser)
	//TODO: Add single device API
	//can create a new "blank" device
	r.GET("/v1/projects/ribbit-test-569244/devices", handleGetRequest)
	//upload device credentials to golioth
	r.GET("/v1/projects/ribbit-test-569244/devices/64194746a946a2ad67aba7ad", handleGetRequest)
	//get the user device identity and PSK
	r.GET("/v1/projects/ribbit-test-569244/credentials", handleGetRequest)

	//TODO: generate credentials for the user

	return r
}
