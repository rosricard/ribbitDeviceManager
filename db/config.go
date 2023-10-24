package db

import "gorm.io/gorm"

type ConfigDB struct {
	apiKey    string `gorm:"column:api_key;primary_key"`
	projectID string `gorm:"column:project_id"`
}

type Config struct {
	APIKey    string
	ProjectID string
}

type ConfigRepo struct {
	db *gorm.DB
}

// TableName sets the table name for the ConfigDB model
func (ConfigDB) TableName() string {
	return "config"
}

// NewConfigRepo initializes a new instance of the [ConfigRepo] type
func NewConfigRepo(db *gorm.DB) *ConfigRepo {
	return &ConfigRepo{db}
}
