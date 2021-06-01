package service

import (
	"github.com/southern-martin/oauth-api/src/domain/user"
	"github.com/southern-martin/oauth-api/src/utils/crypto_util"
	"github.com/southern-martin/user-api/domain/user"
	"github.com/southern-martin/user-api/util/crypto_util"
	"github.com/southern-martin/user-api/util/date_util"
	"github.com/southern-martin/util-go/rest_error"
)

var (
	UserService userServiceInterface = &userService{}
)

type userService struct{}

type userServiceInterface interface {
	GetUser(int64) (*user.User, rest_error.RestErr)
	CreateUser(user.User) (*user.User, rest_error.RestErr)
	UpdateUser(bool, user.User) (*user.User, rest_error.RestErr)
	DeleteUser(int64) rest_error.RestErr
	SearchUser(string) (user.Users, rest_error.RestErr)
	LoginUser(user.LoginRequest) (*user.User, rest_error.RestErr)
}

func (s *userService) GetUser(userId int64) (*user.User, rest_error.RestErr) {
	dao := &user.User{Id: userId}
	if err := dao.Get(); err != nil {
		return nil, err
	}
	return dao, nil
}

func (s *userService) CreateUser(user user.User) (*user.User, rest_error.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Status = user.StatusActive
	user.DateCreated = date_util.GetNowDBFormat()
	user.Password = crypto_util.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) UpdateUser(isPartial bool, user user.User) (*user.User, rest_error.RestErr) {
	current := &user.User{Id: user.Id}
	if err := current.Get(); err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}

		if user.LastName != "" {
			current.LastName = user.LastName
		}

		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}

	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *userService) DeleteUser(userId int64) rest_error.RestErr {
	dao := &user.User{Id: userId}
	return dao.Delete()
}

func (s *userService) SearchUser(status string) (user.Users, rest_error.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

func (s *userService) LoginUser(request user.LoginRequest) (*user.User, rest_error.RestErr) {
	dao := &user.User{
		Email:    request.Email,
		Password: crypto_util.GetMd5(request.Password),
	}
	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}
	return dao, nil
}
