package db

import (
	"time"

	"gorm.io/gorm"
)

type DeviceDB struct {
	GoliothPSKID  string `gorm:"column:golioth_psk_id;primary_key"`
	GoliothPSK    string `gorm:"column:golioth_psk"`
	UserID        string `gorm:"column:user_id"`
	UserGivenName string `gorm:"column:user_given_name"`
	ProjectID     string `gorm:"column:project_id"`
	CreatedAt     time.Time
}

type DeviceRepo struct {
	db *gorm.DB
}

type Device struct {
	GoliothPSKID  string
	GoliothPSK    string
	UserID        string
	UserGivenName string
	ProjectID     string
	CreatedAt     string
}

// TableName sets the table name for the DeviceDB model
func (DeviceDB) TableName() string {
	return "devices"
}

// NewDeviceRepo initializes a new instance of the [UserRepo] type
func NewDeviceRepo(db *gorm.DB) *DeviceRepo {
	return &DeviceRepo{db}
}

// CreateDevice will add a single new device to database

// GetAllDevices retrieves a list of all devices and their info from the database

// DeletedDevice removed a device from the db
