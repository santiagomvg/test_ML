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
const countryInfoCacheKey = "cache:country:%s" //%s=country code
const currencyInfoCacheKey = "cache:currency:global"

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
		abort(c, 500, err)
		return
	}

	cinfo, err := getCountryInfo(ccode)
	if err != nil {
		abort(c, 500, err)
		return
	}

	usdValue, err := getCountryUSDValue(cinfo)
	if err != nil {
		abort(c, 500, err)
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

	s := DB.Session()
	defer s.Close()

	//intento levantar informacion del pais cacheada
	var cinfo CountryInfo
	key := fmt.Sprintf(countryInfoCacheKey, countryCode)
	if err := s.ReadJson(key, &cinfo); err != nil {
		log.Printf("cache missed for %s: %v", countryCode, err)
	}
	s.Close()

	var err error
	if len(cinfo.Alpha3Code) == 0 {

		//no tengo countryInfo, lo pido remotamente
		url := fmt.Sprintf("https://restcountries.eu/rest/v2/alpha/%s", countryCode)
		err = Net.Call(http.MethodGet, url, &cinfo)

		//cacheamos lo obtenido por la api pero en background
		go func(cinfo CountryInfo) {

			defer func() {
				if e := recover(); e != nil {
					log.Printf("countryInfo cache insert error %v", e)
				}
			}()

			s := DB.Session()
			defer s.Close()
			if err := s.StoreJson(key, cinfo); err != nil {
				panic(err) //lo atrapa el recover
			}
			//cache por pais valido durante un dia
			_ = s.ExpiresAt(key, time.Now().Add(24*time.Hour))
		}(cinfo)
	}
	return &cinfo, err
}

func getCountryUSDValue(cinfo *CountryInfo) (float64, error) {

	s := DB.Session()
	defer s.Close()

	//intento levantar informacion de cotizaciones cacheadas
	var currency Currency
	if err := s.ReadJson(currencyInfoCacheKey, &currency); err != nil {
		log.Printf("cache missed for currencies: %v", err)
	}
	s.Close()

	if len(currency.Base) == 0 {

		//no tengo currency cacheada, lo pido remotamente
		url := "http://data.fixer.io/api/latest?access_key=fea0cfce5557c66f2a198a58103e04c2"
		err := Net.Call(http.MethodGet, url, &currency)
		if err != nil {
			return 0, err
		}

		//cacheamos lo obtenido por la api pero en background
		go func(currency Currency) {

			defer func() {
				if e := recover(); e != nil {
					log.Printf("currency cache insert error %v", e)
				}
			}()

			s := DB.Session()
			defer s.Close()
			if err := s.StoreJson(currencyInfoCacheKey, &currency); err != nil {
				panic(err) //lo atrapa el recover
			}

			//cache de cotizaciones dura media hora
			_ = s.ExpiresAt(currencyInfoCacheKey, time.Now().Add(30*time.Minute))
		}(currency)

	}

	//esta api en su version gratuita siempre devuelve cotizaciones con base en euros. Convierto a USD de ser necesario
	usd, exists := currency.Rates["USD"]
	if !exists {
		return 0, nil
	}

	localCurrency := cinfo.Currencies[0].Code
	if localCurrency == "EUR" {
		return usd, nil

	} else if localCurrency == "USD" {
		return 1, nil

	} else {

		local, exists := currency.Rates[localCurrency]
		if !exists {
			return 0, nil
		}
		usdBasedValue := math.Round(((1-math.Abs(1-usd))*local)*100) / 100
		return usdBasedValue, nil
	}

}

func abort(c *gin.Context, status int, e error) {

	type jsonErr struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}
	c.AbortWithStatusJSON(status, &jsonErr{
		Status: status,
		Error:  e.Error(),
	})
}
