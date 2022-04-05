package config

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Username string
	Password string
	Host     string
	DBName   string
}

func Conn() *gorm.DB {
	var cred Config

	cred.Username = os.Getenv("DB_USER")
	cred.Password = os.Getenv("DB_PASS")
	cred.Host = os.Getenv("DB_HOST")
	cred.DBName = os.Getenv("DB_NAME")

	var dns = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", cred.Username, cred.Password, cred.Host, cred.DBName)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		PrepareStmt: true,
	})

	FailOnError(err, 36, "config/database.go")

	return db
}
