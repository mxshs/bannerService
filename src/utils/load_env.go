package utils

import (
	"fmt"
	"os"
	"strconv"
)

func GetPgConnString() string {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		panic("cannot get db hostname from env")
	}

	port, ok := os.LookupEnv("DB_PORT")
	if !ok {
		panic("cannot get db port from env")
	}

	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		panic("cannot get db user from env")
	}

	pass, ok := os.LookupEnv("DB_PASS")
	if !ok {
		panic("cannot get db user password from env")
	}

	dbname, ok := os.LookupEnv("DB_NAME")
	if !ok {
		panic("cannot get db name from env")
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname,
	)
}

func GetRedisConnectionVars() (string, string, int) {
	host, ok := os.LookupEnv("REDIS_HOST")
	if !ok {
		panic("cannot get redis hostname from env")
	}

	port, ok := os.LookupEnv("REDIS_PORT")
	if !ok {
		panic("cannot get redis port from env")
	}

    password, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		panic("cannot get redis password from env")
	}

    rdb, ok := os.LookupEnv("REDIS_DB")
	if !ok {
		panic("cannot get redis db from env")
	}

    db, err := strconv.Atoi(rdb)
    if err != nil {
		panic("failed to convert REDIS_DB env variable to int")
    }

    return fmt.Sprintf("%s:%s", host, port), password, db
}
