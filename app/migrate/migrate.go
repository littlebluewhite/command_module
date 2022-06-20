package migrate

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"new_command/app/database"
)

func NewMigrateDB() (*sql.DB, error) {
	cfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", database.DBConfig.Host, database.DBConfig.Port),
		DBName:               database.DBConfig.DB,
		User:                 database.DBConfig.User,
		Passwd:               database.DBConfig.Password,
		AllowNativePasswords: false,
		ParseTime:            true,
	}

	return sql.Open("mysql", cfg.FormatDSN())
}
