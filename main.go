package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

type config struct {
	Redis struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"redis"`
}

func main() {

	cfg := readConfigFile()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	router.GET("/api/ipinfo", handleIPInfo)
	router.GET("/api/stats/nearest", handleStatsNearest)
	router.GET("/api/stats/farthest", handleStatsFarthest)
	router.GET("/api/stats/avg", handleStatsAVG)

	if err := DB.Init(cfg.Redis.Host, cfg.Redis.Port); err != nil {
		log.Fatalln(err)
	}

	log.Println("Service running")
	log.Println("Open browser in http://localhost:5000")
	if err := router.Run(":5000"); err != nil {
		log.Fatalln(err)
	}
}

func readConfigFile() *config {

	var cfg config
	cfgFile, err := os.Open("config.json")

	if err != nil {
		log.Fatalln(err)
	}
	if err := json.NewDecoder(cfgFile).Decode(&cfg); err != nil {
		log.Fatalln(err)
	}
	return &cfg
}
