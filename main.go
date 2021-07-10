package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/cecilozaur/tnemesum/app"
	"github.com/cecilozaur/tnemesum/conf"
	"github.com/cecilozaur/tnemesum/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	configFile string
	listenPort int
)

func main() {
	flag.StringVar(&configFile, "config", "config.toml", "Config file to use")
	flag.IntVar(&listenPort, "port", 8000, "App listen port")
	flag.Parse()

	cfg := conf.Config{}
	if _, err := toml.DecodeFile(configFile, &cfg); err != nil {
		log.Fatalf("failed to load specified config file")
	}

	muse := app.New(cfg)
	muse.Run()
	muse.SetHealthy()

	r := gin.Default()

	// use middlewares
	r.Use(middleware.RateControl(cfg))

	r.GET("/ping", func(c *gin.Context) {
		if muse.Healthy() {
			c.JSON(http.StatusOK, "OK")
		} else {
			c.JSON(http.StatusBadGateway, "KO")
		}
	})

	api := r.Group("/api")
	{
		api.GET("/cities", muse.GetAllCities)
		api.GET("/cities/:cityId", muse.GetCity)
		api.GET("/forecast/:cityId/", muse.GetCityForecast)
		api.GET("/forecast/:cityId/:day", muse.GetCityForecast)
		api.POST("/forecast/:cityId", muse.UpdateCityForecast)
	}

	log.Fatal(r.Run(fmt.Sprintf(":%d", listenPort)))
}
