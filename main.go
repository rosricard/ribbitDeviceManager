package main

import (
	"log"
	"net/http"

	"github.com/rosricard/userAccess/handlers"

	"github.com/go-chi/chi"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Database func() *gorm.DB
	database *gorm.DB

	Logger func() *log.Logger
	logger *log.Logger
}

type Repository struct {
	users *UserRepo
}

type user struct {
	ID               string `gorm:"column:id;primary_key"`
	name             string `gorm:"column:name"`
	password         string `gorm:"column:password"`
	sensorPrivateKey string `gorm:"column:private_key"`
}

type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo initializes a new instance of the [UserRepo] type
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

// Automigrate automatically updates any configured database table if a mis-match in config is detected
func Automigrate(db *gorm.DB) error {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:password@tcp(127.0.0.1:3306)/ribbit?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&user{})
	return err
}

// NewRepository initializes a new instance of the [Repository] type
//
// [Automigrate] will be called to update the db as part of initialization
func NewRepository(db *gorm.DB) (*Repository, error) {
	if err := Automigrate(db); err != nil {
		return nil, err
	}
	return &Repository{
		users: NewUserRepo(db),
	}, nil
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
