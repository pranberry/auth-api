package config

import "os"

var (
	User     = os.Getenv("DB_USER")
	Password = os.Getenv("DB_PASSWORD")
	DbName   = os.Getenv("DB_NAME")
	Host     = os.Getenv("DB_HOST")
)
