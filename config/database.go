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
	//TODO: godotenv disable
	// err := godotenv.Load()
	// FailOnError(err, 23, "database.go")

	// cred.Username = "u1656216_dk_project_admin"
	// cred.Password = "dk_project_admin_2022"
	// cred.Host = "srv143.niagahoster.com"
	// cred.DBName = "u1656216_dk_database_project"

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
