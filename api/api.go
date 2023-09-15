package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rosricard/userAccess/db"
)

const (
	projectID     = "ribbit-test-569244"
	baseURL       = "https://api.golioth.io"
	apiKey        = "R7aJE5qW4DNHJTgy9JpbmZYrFXnRTY8S"
	getAllDevices = "https://api.golioth.io/v1/projects/ribbit-test-569244/devices/64194746a946a2ad67aba7ad/credentials"
	postDevice    = "https://api.golioth.io/v1/projects/ribbit-test-569244/devices"
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

// goliothGetRequest handles GET requests to the golioth external APIs
func goliothGetRequest(c *gin.Context) (error, *DeviceList) {

	// Create a new HTTP request
	req, err := http.NewRequest("GET", getAllDevices, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return err, nil
	}

	// Add headers to the HTTP request
	req.Header.Set("X-API-Key", apiKey)

	// Make the API call
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
		return err, nil
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return err, nil
	}

	var devices DeviceList

	err1 := json.Unmarshal(body, &devices)
	if err1 != nil {
		return err1, nil
	}

	// return the response to the client
	c.JSON(http.StatusOK, gin.H{"message": resp.Status, "data": devices})
	return nil, &devices
}

func createDevice(c *gin.Context) {

	type goliothDevice struct {
		ProjectID   string   `json:"projectId"`
		Name        string   `json:"name"`
		HardwareIds []string `json:"hardwareIds"`
		TagIds      []string `json:"tagIds"`
		BlueprintId string   `json:"blueprintId"`
	}

	device := goliothDevice{
		ProjectID:   projectID,
		Name:        "Test",
		HardwareIds: []string{"123456789"},
		TagIds:      []string{"string"},
		BlueprintId: "string",
	}

	body, err := json.Marshal(device)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	// Create the API endpoint URL
	url := fmt.Sprintf("%s/v1/projects/%s/devices", baseURL, device.ProjectID)

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	// Add headers to the HTTP request
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("Failed to create device, status code: %d, response: %s", resp.StatusCode, respBody)
	}

	fmt.Println("Successfully created device")

	// return the response to the client
	c.JSON(http.StatusOK, gin.H{"message": resp.Status})
}

//TODO: setup config files with projectID
// user logs in
// add device to table

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/createusers", CreateUser)
	//TODO: change this to GetUser
	r.GET("/getusers", GetAllUsers)
	r.DELETE("/users/:email", DeleteUser)
	r.POST("/createDevice", createDevice)
	//get devices from golioth
	//r.GET("/devices", getAllDevices)
	//get the user device identity and PSK
	//r.GET("/v1/projects/ribbit-test-569244/credentials", getAllDevices)
	// get api keys
	//https://api.golioth.io/v1/projects/ribbit-test-569244/apikeys

	//TODO: generate credentials for the user

	return r
}
