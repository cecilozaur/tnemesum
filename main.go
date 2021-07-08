package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/cecilozaur/tnemesum/app"
	"github.com/cecilozaur/tnemesum/conf"
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
	r.GET("/ping", func(c *gin.Context) {
		if muse.Healthy() {
			c.JSON(http.StatusOK, gin.H{
				"status": "OK",
			})
		} else {
			c.JSON(http.StatusBadGateway, gin.H{
				"status": "KO",
			})
		}
	})

	r.GET("/forecast", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"items": muse.GetAllCities(),
		})
	})

	log.Fatal(r.Run(fmt.Sprintf(":%d", listenPort)))
}
