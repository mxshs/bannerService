package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"mxshs/bannerService/src/adapters/pg"
	"mxshs/bannerService/src/utils"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Printf("%v is not a valid set of arguments\n", os.Args)
        return
    }

    to_gen := os.Args[1]

    num, err := strconv.Atoi(to_gen)
    if err != nil {
        fmt.Printf("%s is not a valid option for the quantity of generated items\n", to_gen)
        return
    }

    if err := godotenv.Load(".env.dev"); err != nil {
        fmt.Println(err.Error())
        return
    }

    db := pg.NewDB(utils.GetPgConnString())

    tag, feature := 1, 1

    for range num {
        tags := make([]int, 0, 5)
        for range 5 {
            tags = append(tags, tag)
            tag++
        }

        randomBytes := map[string]string {
            "field1": genRandomString(255),
            "field2": genRandomString(255),
        }

        content, err := json.Marshal(randomBytes)
        if err != nil {
            fmt.Println(err)
            return
        }

        _, err = db.CreateBanner(tags, feature, content, true)
        if err != nil {
            fmt.Println(err)
            return
        }

        feature++
    }
}

func genRandomString(sz int) string {
    buf := make([]byte, sz)

    for idx := range sz {
        buf[idx] = byte(rand.Intn(26) + 97)
    }

    return string(buf)
}
