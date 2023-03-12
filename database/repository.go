package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DBCon *gorm.DB
)

type Repository struct {
	users *UserRepo
}

type user struct {
	ID               string `gorm:"column:id;primary_key"`
	name             string `gorm:"column:name"`
	password         string `gorm:"column:password"`
	sensorPrivateKey string `gorm:"column:private_key"`
}

func InitDB() {
	var err error

	//connect to mysql database
	dsn := "root:password@tcp(127.0.0.1:3306)/ribbit?charset=utf8mb4&parseTime=True&loc=Local"
	//DBCon, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	DBCon, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	//where myhost is port is the port postgres is running on
	//user is your postgres use name
	//password is your postgres password
	if err != nil {
		panic("failed to connect database")
	}
	print(DBCon)

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
