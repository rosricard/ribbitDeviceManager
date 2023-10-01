package db

import (
	"time"

	"gorm.io/gorm"
)

type DeviceDB struct {
	deviceID   string `gorm:"column:device_id;primary_key"`
	deviceName string `gorm:"column:device_name"`
	devicePSK  string `gorm:"column:device_psk"`
	UserID     string `gorm:"column:user_id"`
	UserName   string `gorm:"column:user_given_name"`
	ProjectID  string `gorm:"column:project_id"`
	CreatedAt  time.Time
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
func (repo *DeviceRepo) CreateDevice(device DeviceDB) error {
	return repo.db.Create(&device).Error
}

// GetAllDevices retrieves a list of all devices and their info from the database

// DeletedDevice removed a device from the db
