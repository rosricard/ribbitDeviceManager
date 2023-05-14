// package db provides functionality for interacting with a relational database
package db

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	users *UserRepo
	//TODO: add in zap logger
}

// AutoMigrate will automatically migrate the database and correcrt schema errors on startup
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&UserDB{})
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &Repository{
		users: NewUserRepo(db),
	}, nil

}

func (r *Repository) Users() *UserRepo {
	return r.users
}
