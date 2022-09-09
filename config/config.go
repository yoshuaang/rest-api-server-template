package config

import "fmt"

var (
	ENV_TEST = "dev" // dev or prod
	// Default Database (MySQL) --> golang_rest_api
	DB_TEST_HOSTNAME = "localhost"
	DB_TEST_NAME     = "golang_rest_api"
	DB_TEST_USERNAME = "root"
	DB_TEST_PASSWORD = ""
	DSN_TEST         = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", DB_TEST_USERNAME, DB_TEST_PASSWORD, DB_TEST_HOSTNAME, DB_TEST_NAME)
)
