// package to store all of the applicaiton data models
package model

import "github.com/gofrs/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:name`
	Password string    `json:password`
}
