package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// online mode

type ResponseB struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data DataB  `json:"data"`
}

type DataB struct {
	IP          string `json:"ip"`
	Dec         string `json:"dec"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Districts   string `json:"districts"`
	IDC         string `json:"idc"`
	ISP         string `json:"isp"`
	Net         string `json:"net"`
	Zipcode     string `json:"zipcode"`
	Areacode    string `json:"areacode"`
	Protocol    string `json:"protocol"`
	Location    string `json:"location"`
	Myip        string `json:"myip"`
	Time        string `json:"time"`
}

func getIpLocationByMir6(ip string) string {
	res, err := http.Get("https://api.mir6.com/api/ip?type=json&ip=" + ip)
	if err != nil {
		return "unknown"
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "unknown"
	}

	var response ResponseB
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "unknown"
	}

	if response.Code != 200 {
		return "unknown"
	}

	country := If(response.Data.Country == "中国" || response.Data.CountryCode == "CN", "中国", response.Data.Country)
	country = If(country != "", country, "0")
	province := If(response.Data.Province != "", response.Data.Province, "0")
	city := If(response.Data.City != "", response.Data.City, "0")
	isp := If(response.Data.ISP != "", response.Data.ISP, "0")
	countryCode := If(response.Data.CountryCode != "", response.Data.CountryCode, "0")

	return fmt.Sprintf("%s|%s|%s|%s|%s|mir6-api|online", country, province, city, isp, countryCode)
}
