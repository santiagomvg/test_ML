package main

/*
Los siguientes structs fueron generados automaticamente por json->GO
Comente varios campos que no necesito para que se pierda tiempo innecesario deserializandolos.
*/

//API: https://api.ip2country.info/ip?186.130.107.97
type IPInfo struct {
	//	CountryCode  string `json:"countryCode"`
	CountryCode3 string `json:"countryCode3"`
	//	CountryName  string `json:"countryName"`
	//	CountryEmoji string `json:"countryEmoji"`
}

//https://restcountries.eu/rest/v2/alpha/ARG
type CountryInfo struct {
	Name string `json:"name"`
	//	TopLevelDomain []string        `json:"topLevelDomain"`
	Alpha2Code string `json:"alpha2Code"`
	Alpha3Code string `json:"alpha3Code"`
	//	CallingCodes   []string        `json:"callingCodes"`
	//	Capital        string          `json:"capital"`
	//	AltSpellings   []string        `json:"altSpellings"`
	//	Region         string          `json:"region"`
	//	Subregion      string          `json:"subregion"`
	//	Population     int             `json:"population"`
	Latlng []float64 `json:"latlng"`
	//	Demonym        string          `json:"demonym"`
	//	Area           float64         `json:"area"`
	//	Gini           float64         `json:"gini"`
	Timezones []string `json:"timezones"`
	//	Borders        []string        `json:"borders"`
	NativeName string `json:"nativeName"`
	//	NumericCode    string          `json:"numericCode"`
	Currencies []Currencies `json:"currencies"`
	Languages  []Languages  `json:"languages"`
	//	Translations   Translations    `json:"translations"`
	//	Flag           string          `json:"flag"`
	//	RegionalBlocs  []RegionalBlocs `json:"regionalBlocs"`
	//	Cioc           string          `json:"cioc"`
}
type Currencies struct {
	Code string `json:"code"`
	Name string `json:"name"`
	//	Symbol string `json:"symbol"`
}
type Languages struct {
	Iso6391    string `json:"iso639_1"`
	Iso6392    string `json:"iso639_2"`
	Name       string `json:"name"`
	NativeName string `json:"nativeName"`
}

/*
type Translations struct {
	De string `json:"de"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	It string `json:"it"`
	Br string `json:"br"`
	Pt string `json:"pt"`
	Nl string `json:"nl"`
	Hr string `json:"hr"`
	Fa string `json:"fa"`
}
type RegionalBlocs struct {
	Acronym       string   `json:"acronym"`
	Name          string   `json:"name"`
	OtherAcronyms []string `json:"otherAcronyms"`
	OtherNames    []string `json:"otherNames"`
}
*/

//fixer (cotizaciones) key = fea0cfce5557c66f2a198a58103e04c2
type Currency struct {
	Success   bool               `json:"success"`
	Timestamp int                `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float64 `json:"rates"`
}
