package main

import (
	"log"
	"net/http"

	"github.com/rosricard/userAccess/api"
	"github.com/rosricard/userAccess/db"
)

func main() {
	db.ConnectDatabase()

	r := api.SetupRouter()

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
