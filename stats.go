package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const nearestCountryKey = "stats:country:nearest"
const farthestCountryKey = "stats:country:farthest"

type countryStat struct {
	Name     string  `json:"name"`
	Distance float64 `json:"distance"`
}

func handleStatsNearest(c *gin.Context) {

	stat, err := getCountryDistance(nearestCountryKey)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, stat)
}

func handleStatsFarthest(c *gin.Context) {
	stat, err := getCountryDistance(farthestCountryKey)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, stat)
}

func handleStatsAVG(c *gin.Context) {
	s := db.Session()
	defer s.Close()
}

func getCountryDistance(key string) (*countryStat, error) {

	s := db.Session()
	defer s.Close()

	var stat countryStat
	if err := s.HGET(key, "country,distance", &stat); err != nil {
		return nil, err
	}
	return &stat, nil
}

func updateStatsForCountry(distance float64, cinfo *CountryInfo) {

	stat, err := getCountryDistance(nearestCountryKey)
	if err == nil {
		if distance < stat.Distance {
			stat.Name = cinfo.Name
			stat.Distance = distance
			s := db.Session()
			defer s.Close()
			s.HSET(nearestCountryKey, "country,distance", &stat)
		}
	}

	stat, err = getCountryDistance(farthestCountryKey)
	if err == nil {
		if distance > stat.Distance {
			stat.Name = cinfo.Name
			stat.Distance = distance
			s := db.Session()
			defer s.Close()
			s.HSET(farthestCountryKey, "country,distance", &stat)
		}
	}
}
