package main

import (
	"log"
	"net/http"

	"github.com/rosricard/userAccess/handlers"

	"github.com/go-chi/chi"

	"gorm.io/gorm"
)

type Config struct {
	Database func() *gorm.DB
	database *gorm.DB

	Logger func() *log.Logger
	logger *log.Logger
}

func main() {
	// config := new(Config)
	// newRepo, err := NewRepository(config.Database())
	// if err != nil {
	// 	//Config.logger.Fatal(err)
	// 	print(err)
	// }
	// print(newRepo)

	// run server locally
	router := chi.NewRouter()
	router.Get("/api/jobs", handlers.GetJobs)
	//run it on port 8080
	err := http.ListenAndServe("0.0.0.0:8080", router)
	if err != nil {
		log.Fatal(err)
	}

}
