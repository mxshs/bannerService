package main

import (
	"fmt"
	"mxshs/bannerService/src/adapters/pg"
	"mxshs/bannerService/src/adapters/redis"
	"mxshs/bannerService/src/handlers"
	"mxshs/bannerService/src/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

    db := pg.NewDB(utils.GetPgConnString())
    bc := redis.NewRedisCache(utils.GetRedisConnectionVars())

    signingKey := os.Getenv("SIGNING_KEY")

    srv := handlers.NewBannerServer(db, []byte(signingKey), db, bc)

    gin.SetMode(gin.ReleaseMode)

    srv.Run()
}
