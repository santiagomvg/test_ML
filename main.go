package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	router.GET("/api/ipinfo", handleIPInfo)
	router.GET("/api/stats/nearest", handleStatsNearest)
	router.GET("/api/stats/farthest", handleStatsFarthest)
	router.GET("/api/stats/avg", handleStatsAVG)

	if err := DB.Init("localhost", 6379); err != nil {
		panic(err)
	}

	log.Println("Service running")
	log.Println("Open browser in http://localhost:5000")
	if err := router.Run(":5000"); err != nil {
		log.Fatalln(err)
	}
}
