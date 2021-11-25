package main

import (
	"log"
	"time"

	"github.com/cmxjs/file-server/config"
	"github.com/cmxjs/file-server/dao"
	"github.com/cmxjs/file-server/models"
	"github.com/cmxjs/file-server/router"
)

func main() {
	// init sqlite database
	if err := dao.InitSqlite(); err != nil {
		log.Fatalln(err)
	}

	// sqlite autoMigrate
	if err := models.InitModel(); err != nil {
		log.Fatalln(err)
	}
	// init redis database
	if err := dao.InitRedis(config.RedisAddr, config.RedisDB); err != nil {
		log.Println(err)
	}

	// delete expired file
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(60))
		for {
			<-ticker.C
			models.DeleteExpiredFile()
		}
	}()

	// gin
	r := router.InitRouter()
	r.Run(":" + config.Port)
}
