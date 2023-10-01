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
)

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
	// Create the API endpoint URL
	url := fmt.Sprintf("%s/v1/projects/%s/tags", baseURL, projectID)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
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
