package database

import "gorm.io/gorm"

var (
	DBCon *gorm.DB
)

func InitDB() {
	var err error

	//connect to postgres database
	//DBCon, err = gorm.Open("postgres", "host=localhost port=myport user=gorm dbname=delivery password=mypassword")

	//connect to mysql database

	//where myhost is port is the port postgres is running on
	//user is your postgres use name
	//password is your postgres password
	if err != nil {

		panic("failed to connect database")
	}

}
