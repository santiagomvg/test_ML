package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	//sort stats:requests:countries by nosort get # get stats:request:country:*
	values, err := redis.Values(s.Raw("SORT", countriesRequestsSetKey,
		"BY", "nosort",
		"GET", "#",
		"GET", "stats:request:country:*"))

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var totalRequestsCount int
	var partialAvg float64
	for n := 0; n < len(values); n += 2 {

		//values[n] = country CODE
		//values[n+1] = requests made

		countryCode := string(values[n].([]uint8))
		cinfo, _ := getCountryInfo(countryCode)
		if len(cinfo.Alpha3Code) == 0 {
			continue
		}

		dist := distance(buenosAiresLat, buenosAiresLng, cinfo.Latlng[0], cinfo.Latlng[1], "K")
		countryRequests, _ := strconv.Atoi(string(values[n+1].([]uint8)))
		partialAvg += dist * float64(countryRequests)
		totalRequestsCount += countryRequests
	}
	averageDistance := partialAvg / float64(totalRequestsCount)

	type ret struct {
		AVG           float64 `json:"avg"`
		TotalRequests int     `json:"totalRequests"`
	}
	c.JSON(200, ret{AVG: averageDistance, TotalRequests: totalRequestsCount})
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
