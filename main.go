package main

import (
	"book-service/config"
	"book-service/src/api"
	"book-service/src/dao"
	"github.com/astaxie/beego"
	"log"
)

func main() {
	config.Init()

	database, err := config.Database()
	if err != nil {
		log.Fatalf("failed to get database configuration: %v", err)
	}
	if err := dao.InitDatabase(database); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	registerRoutes()
	beego.Run(config.GetAppPort())
}

func registerRoutes() {
	beego.Router("/healthcheck", &api.HomeAPI{}, "get:Get")
	beego.Router("/:id([0-9]+", &api.BookAPI{}, "get:Get;delete:Delete;put:Put")
	beego.Router("/all", &api.BookAPI{}, "get:List;post:Post")
}

func hostname() string {
	//hn, err := os.Hostname()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	return config.GetUserHostName()
}
