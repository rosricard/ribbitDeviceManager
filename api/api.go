package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"github.com/rosricard/userAccess/db"
)

const (
	projectID     = "ribbit-test-569244"
	baseURL       = "https://api.golioth.io"
	apiKey        = "R7aJE5qW4DNHJTgy9JpbmZYrFXnRTY8S"
	getAllDevices = "https://api.golioth.io/v1/projects/ribbit-test-569244/devices/64194746a946a2ad67aba7ad/credentials"
	postDevice    = "https://api.golioth.io/v1/projects/ribbit-test-569244/devices"
	tagIds        = "647d5ce530e7d8943a41f874"
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
		ProjectID string   `json:"projectId"` // TODO: query the project id from the golioth API
		Name      string   `json:"name"`
		DeviceIds []string `json:"deviceIds"` // uuid
		TagIds    []string `json:"tagIds"`    // TODO: query the tag id from the golioth API
	}

	// generate device name
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	name := nameGenerator.Generate()

	// generate device id
	did := uuid.New().String()

	device := goliothDevice{
		ProjectID: projectID,
		Name:      name,
		DeviceIds: []string{did},
		TagIds:    []string{tagIds},
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
	// req.Header.Set("Content-Type", "application/json")

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

	// add device to db

	// return the response to the client
	c.JSON(http.StatusOK, gin.H{"message": resp.Status})
}

func getTags(c *gin.Context) {
	getTagsAPI := "https://api.golioth.io/v1/projects/" + projectID + "/tags"

	// Create a new HTTP request
	req, err := http.NewRequest("GET", getTagsAPI, nil)
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

	// return the response to the client
	c.JSON(http.StatusOK, gin.H{"message": resp.Status, "data": string(body)})
}

//TODO: setup config files with projectID, tagIds, APIkey, etc
// user logs in
// add device to table

//TODO: on app startup, run a check against the golioth API to get all devices and compare against the database

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/createusers", CreateUser)
	//TODO: change this to GetUser
	r.GET("/getusers", GetAllUsers)
	r.GET("/getTags", getTags)
	r.DELETE("/users/:email", DeleteUser)
	r.POST("/createDevice", createDevice)
	//r.GET("/devices", getAllDevices)
	// get the user device identity and PSK
	// r.GET("/v1/projects/ribbit-test-569244/credentials", getAllDevices)
	// get api keys
	//https://api.golioth.io/v1/projects/ribbit-test-569244/apikeys
	return r
}
