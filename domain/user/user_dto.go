package user

import (
	"github.com/southern-martin/util-go/rest_error"
	"strings"
)

const (
	StatusActive = "active"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

type Users []User

func (user *User) Validate() rest_error.RestErr {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return rest_error.NewBadRequestError("invalid email address")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return rest_error.NewBadRequestError("invalid password")
	}
	return nil
}
