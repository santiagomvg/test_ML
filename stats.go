package main

import (
	"fmt"
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

	stat, err := getCountryDistance(nearestCountryKey, "ARG:BA")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, stat)
}

func handleStatsFarthest(c *gin.Context) {
	stat, err := getCountryDistance(farthestCountryKey, "ARG:BA")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, stat)
}

func handleStatsAVG(c *gin.Context) {
	s := DB.Session()
	defer s.Close()
}

func getCountryDistance(key string, from string) (*countryStat, error) {

	s := DB.Session()
	defer s.Close()

	var stat countryStat
	key = fmt.Sprintf(key+":%s:%s", from)
	if err := s.HGET(key, "country,distance", &stat); err != nil {
		return nil, err
	}
	return &stat, nil
}

func updateStatsForCountry(distance float64, fromGeoLocationCode string, cinfo *CountryInfo) {

	stat, err := getCountryDistance(nearestCountryKey, fromGeoLocationCode)
	if err == nil {
		if distance < stat.Distance {
			stat.Name = cinfo.Name
			stat.Distance = distance
			s := DB.Session()
			defer s.Close()
			s.HSET(nearestCountryKey, "country,distance", &stat)
		}
	}

	stat, err = getCountryDistance(farthestCountryKey, fromGeoLocationCode)
	if err == nil {
		if distance > stat.Distance {
			stat.Name = cinfo.Name
			stat.Distance = distance
			s := DB.Session()
			defer s.Close()
			s.HSET(farthestCountryKey, "country,distance", &stat)
		}
	}
}
