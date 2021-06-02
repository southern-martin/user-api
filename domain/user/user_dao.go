package user

import (
	"errors"
	"fmt"
	"github.com/southern-martin/user-api/datasource/mysql/user_db"
	"github.com/southern-martin/user-api/util/mysql_util"
	"github.com/southern-martin/util-go/logger"
	"github.com/southern-martin/util-go/rest_error"
	"strings"
)

const (
	queryInsertUser             = "INSERT INTO user(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, date_created, status FROM user WHERE id=?;"
	queryUpdateUser             = "UPDATE user SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser             = "DELETE FROM user WHERE id=?;"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM user WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM user WHERE email=? AND password=? AND status=?"
)

func (user *User) Get() rest_error.RestErr {
	stmt, err := user_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_error.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("error when trying to get user by id", getErr)
		return rest_error.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}
	return nil
}

func (user *User) Save() rest_error.RestErr {
	stmt, err := user_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_error.NewInternalServerError("error when tying to save user", errors.New("database error"))
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return rest_error.NewInternalServerError("error when tying to save user", errors.New("database error"))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return rest_error.NewInternalServerError("error when tying to save user", errors.New("database error"))
	}
	user.Id = userId

	return nil
}

func (user *User) Update() rest_error.RestErr {
	stmt, err := user_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_error.NewInternalServerError("error when tying to update user", errors.New("database error"))
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return rest_error.NewInternalServerError("error when tying to update user", errors.New("database error"))
	}
	return nil
}

func (user *User) Delete() rest_error.RestErr {
	stmt, err := user_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_error.NewInternalServerError("error when tying to update user", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return rest_error.NewInternalServerError("error when tying to save user", errors.New("database error"))
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, rest_error.RestErr) {
	stmt, err := user_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user by status statement", err)
		return nil, rest_error.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find user by status", err)
		return nil, rest_error.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when scan user row into user struct", err)
			return nil, rest_error.NewInternalServerError("error when tying to gett user", errors.New("database error"))
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, rest_error.NewNotFoundError(fmt.Sprintf("no user matching status %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() rest_error.RestErr {
	stmt, err := user_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_error.NewInternalServerError("error when tying to find user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_util.ErrorNoRows) {
			return rest_error.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return rest_error.NewInternalServerError("error when tying to find user", errors.New("database error"))
	}
	return nil
}
