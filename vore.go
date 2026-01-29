package main

// online mode

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResponseA struct {
	Code   int     `json:"code"`
	Msg    string  `json:"msg"`
	IPInfo IPInfoA `json:"ipinfo"`
	IPData IPDataA `json:"ipdata"`
	Adcode AdcodeA `json:"adcode"`
	Tips   string  `json:"tips"`
	Time   int64   `json:"time"`
}

type IPInfoA struct {
	Type string `json:"type"`
	Text string `json:"text"`
	CNIP bool   `json:"cnip"`
}

type IPDataA struct {
	Info1 string `json:"info1"`
	Info2 string `json:"info2"`
	Info3 string `json:"info3"`
	ISP   string `json:"isp"`
}

type AdcodeA struct {
	O string      `json:"o"`
	P string      `json:"p"`
	C string      `json:"c"`
	N string      `json:"n"`
	R interface{} `json:"r"`
	A interface{} `json:"a"`
	I bool        `json:"i"`
}

func getIpLocationByVore(ip string) string {
	res, err := http.Get("https://api.vore.top/api/IPdata?ip=" + ip)
	if err != nil {
		return "unknown"
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "unknown"
	}

	var response ResponseA
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "unknown"
	}

	if response.Code != 200 {
		return "unknown"
	}

	var country, subdivision, city, isp string
	if response.IPInfo.CNIP {
		country = "中国"
		subdivision = response.IPData.Info1
		city = response.IPData.Info2
		isp = response.IPData.ISP
	} else {
		country = response.IPData.Info1
		subdivision = response.IPData.Info2
		city = response.IPData.Info3
		isp = response.IPData.ISP
	}

	country = If(country == "", "0", country)
	subdivision = If(subdivision == "", "0", subdivision)
	city = If(city == "", "0", city)
	isp = If(isp == "", "0", isp)

	return fmt.Sprintf("%s|%s|%s|%s|vore-api|online", country, subdivision, city, isp)
}
