package user_db

import (
	"database/sql"
	"fmt"
	"github.com/southern-martin/utils-go/logger"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const (
	mysqlUserUsername = "mysql_users_username"
	mysqlUserPassword = "mysql_users_password"
	mysqlUserHost     = "mysql_users_host"
	mysqlUserSchema   = "mysql_users_schema"
)

var (
	Client *sql.DB

	username = os.Getenv(mysqlUserUsername)
	password = os.Getenv(mysqlUserPassword)
	host     = os.Getenv(mysqlUserHost)
	schema   = os.Getenv(mysqlUserSchema)
)

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username, password, host, schema,
	)
	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		panic(err)
	}

	mysql.SetLogger(logger.GetLogger())
	log.Println("database successfully configured")
}
