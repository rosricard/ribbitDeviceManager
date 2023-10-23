package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rosricard/ribbitDeviceManager/db"
	"golang.org/x/crypto/bcrypt"
)

//TODO: store config in db a json string, create config file to handle creation and parsing of these details

const (
	projectID = "ribbit-test-569244"
	baseURL   = "https://api.golioth.io"
	apiKey    = "R7aJE5qW4DNHJTgy9JpbmZYrFXnRTY8S"
)

type Device struct {
	ID           string
	Name         string
	PreSharedKey string
	UserID       string
	ProjectID    string
	CreatedAt    time.Time
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Initialize the session store
var store = sessions.NewCookieStore([]byte("some-secret-key"))

func Signup(c *gin.Context) {
	creds := &Credentials{
		Email:    c.Param("email"),
		Password: c.Param("password"),
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	hashedPassword := string(hashedPasswordBytes)

	user := db.User{
		Email:    creds.Email,
		Password: hashedPassword,
	}

	if err := db.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// After successfully creating the user
	session, _ := store.Get(c.Request, "user-session")
	session.Values["email"] = creds.Email
	session.Values["lastActivity"] = time.Now()
	session.Save(c.Request, c.Writer)

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}

func Signin(c *gin.Context) {
	creds := &Credentials{
		Email:    c.Param("email"),
		Password: c.Param("password"),
	}

	// Get the existing entry present in the database for the given username
	user, err := db.GetUserByEmail(creds.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	storedCreds := &Credentials{
		Password: user.Password,
		Email:    user.Email,
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	session, _ := store.Get(c.Request, "user-session")
	session.Values["email"] = creds.Email
	session.Values["lastActivity"] = time.Now()
	session.Save(c.Request, c.Writer)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func activityMiddleware(c *gin.Context) {
	session, _ := store.Get(c.Request, "user-session")
	lastActivity, ok := session.Values["lastActivity"].(time.Time)

	if ok && time.Since(lastActivity) <= 15*time.Minute {
		// Update the last activity time
		session.Values["lastActivity"] = time.Now()
		session.Save(c.Request, c.Writer)
		c.Next()
	} else {
		session.Options.MaxAge = -1 // Invalidate session
		session.Save(c.Request, c.Writer)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
	}
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

// addDeviceToUser adds a device to the active user account
func addDeviceToUser(c *gin.Context) {
	// Access the session
	session, err := store.Get(c.Request, "user-session")
	if err != nil {
		return
	}

	// Retrieve the active user's email from the session
	email, ok := session.Values["email"].(string)
	if !ok {
		return
	}

	// Fetch the user details using the email from the session
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return
	}

	//create device
	d, err := createNewDevice()
	if err != nil {
		return
	}

	//create private key
	psk, err := createPSK(d.DeviceId)
	if err != nil {
		return
	}

	device := db.DeviceDB{
		DeviceID:   d.DeviceId,
		DeviceName: d.Name,
		DevicePSK:  psk.PreSharedKey,
		UserID:     user.ID,
		ProjectID:  d.ProjectID,
		CreatedAt:  psk.CreatedAt,
	}

	err1 := db.CreateDevice(device)
	if err1 != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"deviceID": d.DeviceId, "psk": psk.PreSharedKey, "email": email})

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

	c.JSON(http.StatusOK, gin.H{"deviceID": device.DeviceId, "psk": pskData.Identity})

}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/signin/:email/:password", Signin)
	r.POST("/signup/:email/:password", Signup)
	r.Use(activityMiddleware) // Use the middleware to track active user login status
	r.GET("/getusers", GetAllUsers)
	r.DELETE("/users/:email", DeleteUser)
	r.POST("/createDevice", addDeviceToUser)
	return r
}

// TODO: setup config files with projectID, tagIds, APIkey, etc
// add device to table
// TODO: on app startup, run a check against the golioth API to get all devices and compare against the database
