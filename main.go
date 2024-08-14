package main

import (
	"log"
	"net/http"

	"github.com/igorlopes88/goexpert-ratelimiter/configs"
	"github.com/igorlopes88/goexpert-ratelimiter/infra/middlewares"
	"github.com/igorlopes88/goexpert-ratelimiter/infra/storage"
	"github.com/igorlopes88/goexpert-ratelimiter/infra/webserver"
)

func main() {
	// LOAD .ENV
	conf, err := configs.Load(".")
	if err != nil {
		panic(err)
	}

	// DEFINE STORAGE
	storageAdapter, err := storage.InitRedis(conf.DbAddress, conf.DbPort)
	if err != nil {
		panic(err)
	}

	// CREATE WEBSERVER
	webServer := webserver.New(conf.WebServerPort)
	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`authorized request`))
	}
	webServer.AddHandler("/", mainHandler, "GET")

	// ADD MIDDLEWARE (RATELIMIT)
	settings, err := middlewares.Load(*conf, storageAdapter)
	if err != nil {
		panic(err)
	}
	rateLimiter := middlewares.NewRateLimiter(settings)
	webServer.AddMiddleware(rateLimiter)

	// INIT WEBSERVER
	log.Print("server up in http://localhost:8080")
	webServer.Start()
	log.Print("server down")
}
