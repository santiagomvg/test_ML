package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net"
	"net/http"
)

const buenosAiresLat = 0
const buenosAiresLng = 0

type APIResult struct {
	CountryName string  `json:"countryName"`
	CountryCode string  `json:"countryCode"`
	Distance    float64 `json:"distance"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Currency    string  `json:"currency"`
	USDValue    float64 `json:"usdValue"`
}

func handleIPInfo(c *gin.Context) {

	ip, _ := c.GetQuery("ip")
	ccode, err := getCountryCodeFromIP(ip)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	cinfo, err := getCountryInfo(ccode)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	usdValue, err := getCountryUSDValue(cinfo)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	dist := distance(buenosAiresLat, buenosAiresLng, cinfo.Latlng[0], cinfo.Latlng[1], "K")
	out := APIResult{
		CountryCode: cinfo.Alpha3Code,
		CountryName: cinfo.Name,
		Distance:    dist,
		Latitude:    cinfo.Latlng[0],
		Longitude:   cinfo.Latlng[1],
		Currency:    cinfo.Currencies[0].Code,
		USDValue:    usdValue,
	}

	go updateStatsForCountry(dist, cinfo)
	c.JSON(200, &out)
}

func getCountryCodeFromIP(ipStr string) (string, error) {

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", errors.New("invalid ip address")
	}

	var data IPInfo
	url := fmt.Sprintf("https://api.ip2country.info/ip?%s", ip)
	err := Net.Call(http.MethodGet, url, &data)
	return data.CountryCode3, err
}

func getCountryInfo(countryCode string) (*CountryInfo, error) {

	var data CountryInfo
	url := fmt.Sprintf("https://restcountries.eu/rest/v2/alpha/%s", countryCode)
	err := Net.Call(http.MethodGet, url, &data)
	return &data, err
}

func getCountryUSDValue(cinfo *CountryInfo) (float64, error) {

	var data Currency
	url := "http://data.fixer.io/api/latest?access_key=fea0cfce5557c66f2a198a58103e04c2"
	err := Net.Call(http.MethodGet, url, &data)
	if err != nil {
		return 0, err
	}

	//esta api en su version gratis siempre devuelve cotizaciones con base en euros. Convierto a USD
	local, exists := data.Rates[cinfo.Currencies[0].Code]
	if !exists {
		return 0, nil
	}

	usd, exists := data.Rates["USD"]
	if !exists {
		return 0, nil
	}

	usdBasedValue := (1 - math.Abs(1-usd)) * local
	return usdBasedValue, nil
}
