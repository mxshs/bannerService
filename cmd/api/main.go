package main

import (
	"fmt"
	"mxshs/bannerService/src/adapters/pg"
	"mxshs/bannerService/src/adapters/redis"
	"mxshs/bannerService/src/handlers"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func getPgConnString() string {
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

func getRedisConnectionVars() (string, string, int) {
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


func main() {
    enviroment := "dev"

    if len(os.Args) > 1 {
        enviroment = os.Args[1]
    } else {
        fmt.Println("[INFO] consider providing environment name to load a specific config")
    }

    if err := godotenv.Load(fmt.Sprintf(".env.%s", enviroment)); err != nil {
        panic(err)
    }

    db := pg.NewDB(getPgConnString())
    bc := redis.NewRedisCache(getRedisConnectionVars())

    signingKey := os.Getenv("SIGNING_KEY")

    srv := handlers.NewBannerServer(db, []byte(signingKey), db, bc)

    srv.Run()
}
