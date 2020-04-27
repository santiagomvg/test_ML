package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const nearestCountryKey = "stats:country:nearest:%s"       //nearest to country %s
const farthestCountryKey = "stats:country:farthest:%s"     //farthest to country %s
const countriesRequestsSetKey = "stats:requests:countries" //redis set of all country codes that called
const countryRequestCount = "stats:request:country:%s"     //request count for country %$

type countryStat struct {
	Country  string  `json:"country"`
	Distance float64 `json:"distance"`
}

func handleStatsNearest(c *gin.Context) {

	stat, err := getStoredCountryDistance(nearestCountryKey, "ARG:BA")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, stat)
}

func handleStatsFarthest(c *gin.Context) {
	stat, err := getStoredCountryDistance(farthestCountryKey, "ARG:BA")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, stat)
}

func handleStatsAVG(c *gin.Context) {
	s := DB.Session()
	defer s.Close()

	//TODO:
	//sort stats:requests:countries by nosort get # get stats:request:country:*
}

func getStoredCountryDistance(key string, from string) (*countryStat, error) {

	var stat countryStat
	key = fmt.Sprintf(key, from)

	s := DB.Session()
	defer s.Close()

	if err := s.HGETALL(key, &stat); err != nil {
		return nil, err
	}
	return &stat, nil
}

func updateStatsForCountry(distance float64, fromGeoLocationCode string, cinfo *CountryInfo) {

	//TODO: meter todo en go funcs individuales?
	stat, err := getStoredCountryDistance(nearestCountryKey, fromGeoLocationCode)
	if err == nil {
		if distance < stat.Distance {
			stat.Country = cinfo.Name
			stat.Distance = distance
			s := DB.Session()
			defer s.Close()

			key := fmt.Sprintf(nearestCountryKey, fromGeoLocationCode)
			s.HMSET(key, stat)
		}
	}

	stat, err = getStoredCountryDistance(farthestCountryKey, fromGeoLocationCode)
	if err == nil {
		if distance > stat.Distance {
			stat.Country = cinfo.Name
			stat.Distance = distance
			s := DB.Session()
			defer s.Close()

			key := fmt.Sprintf(farthestCountryKey, fromGeoLocationCode)
			s.HMSET(key, stat)
		}
	}

	s := DB.Session()
	defer s.Close()

	s.SetADD(countriesRequestsSetKey, cinfo.Alpha3Code)
	s.INC(fmt.Sprintf(countryRequestCount, cinfo.Alpha3Code))
}
