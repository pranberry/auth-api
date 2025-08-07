package config

import (
	"time"
)

// var (
// 	User     = os.Getenv("DB_USER")
// 	Password = os.Getenv("DB_PASSWORD")
// 	Dbname   = os.Getenv("DB_NAME")
// 	Host     = os.Getenv("DB_HOST")
// 	TestDb	 = os.Getenv("DB_TEST")
// )

var (
	User     = "token_master"
	Password = "eipu9ahKai2oo9phaib"
	Dbname   = "jwt_users"
	Host     = "localhost"
	TestDb   = "jwt_test"
)

// Token related vars
var JwtTTL = (15 * time.Minute)

const TokenIssuer = "SCDP"
