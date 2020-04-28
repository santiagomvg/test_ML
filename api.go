package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"time"
)

const buenosAiresLat = -34.6131516
const buenosAiresLng = -58.3772316

type APIResult struct {
	CurrentTime  string      `json:"currentTime"`
	NativeName   string      `json:"nativeName"`
	CountryName  string      `json:"countryName"`
	CountryCode  string      `json:"countryCode"`
	Languages    []Languages `json:"languages"`
	Timezones    []string    `json:"timezones"`
	Distance     float64     `json:"distance"`
	DistanceUnit string      `json:"distanceUnit"`
	Latitude     float64     `json:"latitude"`
	Longitude    float64     `json:"longitude"`
	Currency     string      `json:"currency"`
	USDValue     float64     `json:"usdValue"`
}

func handleIPInfo(c *gin.Context) {

	ip, _ := c.GetQuery("ipAddress")
	ccode, err := getCountryCodeFromIP(ip)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
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
		CurrentTime:  time.Now().Format(time.RFC822Z),
		NativeName:   cinfo.NativeName,
		CountryCode:  cinfo.Alpha2Code,
		CountryName:  cinfo.Name,
		Languages:    cinfo.Languages,
		Timezones:    getLocalTimes(cinfo.Timezones),
		Distance:     dist,
		DistanceUnit: "km",
		Latitude:     cinfo.Latlng[0],
		Longitude:    cinfo.Latlng[1],
		Currency:     cinfo.Currencies[0].Code,
		USDValue:     usdValue,
	}

	go updateStatsForCountry(dist, "ARG:BA", cinfo)
	c.JSON(200, &out)
}

func getLocalTimes(tzOffsets []string) (ret []string) {

	ret = make([]string, len(tzOffsets))
	defer func() {
		if e := recover(); e != nil {
			log.Printf("localtime convert error %v", e)
			return
		}
	}()

	utcTime := time.Now().In(time.UTC)
	for i, tzo := range tzOffsets {

		if tzo == "UTC" {
			ret[i] = utcTime.Format(time.Kitchen) + " UTC"

		} else if len(tzo) > 3 && tzo[:3] == "UTC" && len(tzo) > 5 {
			ofsStr := tzo[3:][:3]
			ofs, err := strconv.Atoi(ofsStr)
			if err != nil {
				ret[i] = ""
			} else {
				t := utcTime.Add(time.Duration(ofs) * time.Hour)
				ret[i] = t.Format(time.Kitchen) + " UTC" + ofsStr
			}
		}
	}
	return ret
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

	//esta api en su version gratuita siempre devuelve cotizaciones con base en euros. Convierto a USD de ser necesario
	usd, exists := data.Rates["USD"]
	if !exists {
		return 0, nil
	}

	localCurrency := cinfo.Currencies[0].Code
	if localCurrency == "EUR" {
		return usd, nil

	} else if localCurrency == "USD" {
		return 1, nil

	} else {

		local, exists := data.Rates[localCurrency]
		if !exists {
			return 0, nil
		}
		usdBasedValue := math.Round(((1-math.Abs(1-usd))*local)*100) / 100
		return usdBasedValue, nil
	}
}
