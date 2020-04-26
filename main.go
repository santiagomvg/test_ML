package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "static")

	router.GET("/ping", func(c *gin.Context) {
		c.Writer.Write([]byte("HI!"))
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/api/ipinfo", handleIPInfo)
	router.GET("/api/stats/nearest", handleStatsNearest)
	router.GET("/api/stats/farthest", handleStatsFarthest)
	router.GET("/api/stats/avg", handleStatsAVG)

	log.Println("Service running")
	log.Println("Open browser in http://localhost:5000")
	router.Run(":5000")
}
