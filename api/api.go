package api

import (
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
		ID:        userInput.ID,
		Name:      userInput.Name,
		Email:     userInput.Email,
		Password:  userInput.Password, //TODO: hash password
		ProjectID: userInput.ProjectID,
		DeviceID:  userInput.DeviceID, //TODO: retrieve private key from golioth api instead of user input
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
	userID := c.Param("id")

	if err := db.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// getGoliothPrivateKey retrieves the private key from the golioth api

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/createusers", CreateUser)
	r.GET("/getusers", GetAllUsers)
	r.DELETE("/users/:id", DeleteUser)

	return r
}
