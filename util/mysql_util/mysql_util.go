package mysql_util

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/southern-martin/util-go/rest_error"
	"strings"
)

const (
	ErrorNoRows = "no rows in result set"
)

func ParseError(err error) rest_error.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return rest_error.NewNotFoundError("no record matching given id")
		}
		return rest_error.NewInternalServerError("error parsing database response", err)
	}

	switch sqlErr.Number {
	case 1062:
		return rest_error.NewBadRequestError("invalid data")
	}
	return rest_error.NewInternalServerError("error processing request", errors.New("database error"))
}
