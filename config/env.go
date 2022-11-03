package config

import (
	"os"
)

// TODO: disable when to production
// var (
// 	_ = godotenv.Load()
// )

func GetEnvWA() (string, string) {
	accSID := os.Getenv("T_ACC_SID")
	authToken := os.Getenv("T_AUTH_TOKEN")

	return accSID, authToken
}

func GetEnvDB() Config {
	var username, password, host, name string

	username = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASS")
	host = os.Getenv("DB_HOST")
	name = os.Getenv("DB_NAME")

	return Config{
		Username: username,
		Password: password,
		Host:     host,
		DBName:   name,
	}
}

// first return is USER_KEY, second return is PASS_KEY
func GetEnvZenziva() (string, string) {
	userKey := os.Getenv("ZENZIVA_USER_KEY")
	passKey := os.Getenv("ZENZIVA_PASS_KEY")

	return userKey, passKey
}

func GetEnvJWTKey() string {
	key := os.Getenv("JWT_SECRET")
	return key
}
