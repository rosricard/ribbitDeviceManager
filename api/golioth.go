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

type pskRespData struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Identity     string    `json:"identity"`
	CreatedAt    time.Time `json:"createdAt"`
	PreSharedKey string    `json:"preSharedKey"`
}

type newDevice struct {
	ID          string   `json:"id"`
	hardwareIDs []string `json:"hardwareIds"`
}

// createDevice creates a new device and returns the device id and psk
func createDevice(c *gin.Context) {
	createNewDevice()
}

// createDevice calls the golioth API to create a new device
func createNewDevice() (newDevice, error) {

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

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to create device, status code: %d, response: %s", resp.StatusCode, respBody)
	}

	//unmarshal response into newDevice struct
	var newDevice newDevice
	if err := json.Unmarshal([]byte(respBody), &newDevice); err != nil {
		fmt.Println("Error:", err)
		return newDevice, err
	}

	return newDevice, nil
}

// createPrivateKey creates a private key for the device after the device itself has been created
func createPSK(deviceID string) (pskRespData, error) {
	type goliothPSKreq struct {
		PreSharedKey string `json:"preSharedKey"`
	}

	psk := goliothPSKreq{
		PreSharedKey: "string",
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

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to create device psk, status code: %d, response: %s", resp.StatusCode, respBody)
	}

	//unmarshal response
	var pskData pskRespData
	if err := json.Unmarshal([]byte(respBody), &pskData); err != nil {
		fmt.Println("Error:", err)
		return pskData, err
	}

	return pskData, nil

}
