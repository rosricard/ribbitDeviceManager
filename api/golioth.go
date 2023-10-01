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
	//create private key

	//getUser info
	//combine user info and device info

	// add device to db if success was confirmed

	// return the response to the client
	c.JSON(http.StatusOK, gin.H{"message": resp.Status})
}

// createPrivateKey creates a private key for the device after the device itself has been created
func createPSK(projectID, deviceID string) (string, error) {
	type goliothPSKreq struct {
		ProjectID    string `json:"projectId"`
		preSharedKey string `json:"preSharedKey"`
	}

	psk := goliothPSKreq{
		ProjectID:    projectID,
		preSharedKey: "string",
	}

	body, err := json.Marshal(psk)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	// Create the API endpoint URL
	url := fmt.Sprintf("%s/v1/projects/%s/devices/%s/credentials", baseURL, projectID, deviceID)

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

	//unmarshal response
	type pskRespData struct {
		ID           string    `json:"id"`
		Type         string    `json:"type"`
		Identity     string    `json:"identity"`
		CreatedAt    time.Time `json:"createdAt"`
		PreSharedKey string    `json:"preSharedKey"`
	}

	var pskData pskRespData
	if err := json.Unmarshal([]byte(respBody), &pskData); err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return pskData.PreSharedKey, nil

}

// test createPSK func
func createDevicePrivateKey(c *gin.Context) {
	createPSK(projectID, "6518abaaebe0f4c62ee15eb6")
}
