package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/oschwald/geoip2-golang"
)

//go:embed GeoLite2-City.mmdb
var geoLite2CityFS embed.FS

func getGeoLite2CityMMDB() (*geoip2.Reader, error) {
	geoLite2CityBytes, err := geoLite2CityFS.ReadFile("GeoLite2-City.mmdb")
	if err != nil {
		return nil, err
	}
	return geoip2.FromBytes(geoLite2CityBytes)
}

func getIPLocation(db *geoip2.Reader, ip string) string {
	record, err := db.City(net.ParseIP(ip))
	if err != nil {
		return "unknown"
	}
	country, ok := record.Country.Names["en"]
	if !ok {
		return "unknown"
	}
	if len(record.Subdivisions) == 0 {
		return fmt.Sprintf("[%s]", country)
	}
	subdivision, ok := record.Subdivisions[0].Names["en"]
	if !ok {
		return fmt.Sprintf("[%s]", country)
	}
	city, ok := record.City.Names["en"]
	if !ok {
		return fmt.Sprintf("[%s][%s]", country, subdivision)
	}
	return fmt.Sprintf("[%s][%s][%s]", country, subdivision, city)
}

type IPInfo struct {
	Type string `json:"type"`
	Text string `json:"text"`
	CNIP bool   `json:"cnip"`
}

type IPData struct {
	Info1 string `json:"info1"`
	Info2 string `json:"info2"`
	Info3 string `json:"info3"`
	ISP   string `json:"isp"`
}

type Adcode struct {
	O string      `json:"o"`
	P string      `json:"p"`
	C string      `json:"c"`
	N string      `json:"n"`
	R interface{} `json:"r"`
	A interface{} `json:"a"`
	I bool        `json:"i"`
}

type Response struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	IPInfo IPInfo `json:"ipinfo"`
	IPData IPData `json:"ipdata"`
	Adcode Adcode `json:"adcode"`
	Tips   string `json:"tips"`
	Time   int64  `json:"time"`
}

func getIPLocationFromAPI(ip string) string {
	res, err := http.Get("https://api.vore.top/api/IPdata?ip=" + ip)
	if err != nil {
		return "unknown"
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "unknown"
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "unknown"
	}
	if response.Code != 200 {
		return "unknown"
	}
	var country, subdivision, city, isp string
	if response.IPData.Info1 != "" {
		country = fmt.Sprintf("[%s]", response.IPData.Info1)
	}
	if response.IPData.Info2 != "" {
		subdivision = fmt.Sprintf("[%s]", response.IPData.Info2)
	}
	if response.IPData.Info3 != "" {
		city = fmt.Sprintf("[%s]", response.IPData.Info3)
	}
	if response.IPData.ISP != "" {
		isp = fmt.Sprintf("[%s]", response.IPData.ISP)
	}
	return country + subdivision + city + isp
}

//go:embed nodes.txt
var nodeCIDRsString string

func getNodeCIDRs(nodeCIDRsString string) (string, []*net.IPNet, error) {
	var nodeDate string
	nodeCIDRsSlice := strings.Split(strings.ReplaceAll(nodeCIDRsString, "\r", ""), "\n")
	nodeDate = nodeCIDRsSlice[0]
	_, err := time.Parse("2006-01-02", nodeDate)
	if err != nil {
		return "", nil, err
	}
	nodeCIDRsSlice = nodeCIDRsSlice[1:]
	var nodeCIDRs []*net.IPNet
	for _, nodeCIDR := range nodeCIDRsSlice {
		if nodeCIDR != "" {
			_, nodeCIDR, err := net.ParseCIDR(nodeCIDR)
			if err != nil {
				return "", nil, err
			}
			nodeCIDRs = append(nodeCIDRs, nodeCIDR)
		}
	}
	return nodeDate, nodeCIDRs, nil
}

func printNodeCIDRs(nodeCIDRs []*net.IPNet) {
	fmt.Println("Node IP CIDRs:")
	for i, nodeCIDR := range nodeCIDRs {
		fmt.Printf("%*s) %s\n", 4, strconv.Itoa(i+1), nodeCIDR.String())
	}
}

func initFlag() (string, string, bool) {
	flag.Usage = func() {
		println("Usage:")
		println("  checkKnownsecNodeIP [-a value] [-f value] [-p] [-h]")
		println("  '-a' is necessary argument")
		println("Options:")
		flag.PrintDefaults()
		println("Examples:")
		println("  checkKnownsecNodeIP -a 1.1.1.1")
		println("  checkKnownsecNodeIP -a 1.1.1.1 -f ip.txt")
		println("  checkKnownsecNodeIP -p")
	}
	ipStr := flag.String("a", "", "ip address")
	nodeFile := flag.String("f", "", "node file path")
	printNodes := flag.Bool("p", false, "print node ip list")
	flag.Parse()
	return *ipStr, *nodeFile, *printNodes
}

var version string = "local-build"

func main() {
	nodeDate, nodeCIDRs, err := getNodeCIDRs(nodeCIDRsString)
	if err != nil {
		println("File nodes.txt format error, the first line must be a date (format 'YYYY-mm-dd') and the following lines must be legal CIDRs format")
		return
	}

	db, err := getGeoLite2CityMMDB()
	if err != nil {
		println("GeoLite2 city database error")
		return
	}
	defer db.Close()

	println("checkKnownsecNodeIP " + version + ", A tool to check if an IP is a knownsec node ip")
	println("Node IPs Update At:     " + nodeDate)
	println("IP Database Update At:  " + time.Unix(int64(db.Metadata().BuildEpoch), 0).Format("2006-01-02"))
	println()

	ipStr, nodeFile, printNodes := initFlag()
	if nodeFile != "" {
		nodeCIDRsBytesSlice, err := os.ReadFile(nodeFile)
		if err != nil {
			println("Read file error")
			return
		}
		nodeDate, nodeCIDRs, err = getNodeCIDRs(string(nodeCIDRsBytesSlice))
	}
	if printNodes {
		printNodeCIDRs(nodeCIDRs)
		return
	}
	if ipStr == "" {
		flag.Usage()
		return
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		println("IP address format error")
		return
	}

	println("IP:           " + ipStr)
	println("Location[1]:  " + getIPLocation(db, ipStr))
	println("Location[2]:  " + getIPLocationFromAPI(ipStr))
	for _, nodeCIDR := range nodeCIDRs {
		if nodeCIDR.Contains(ip) {
			println("Node IP:      " + color.New(color.BgGreen, color.Bold).Sprint(" Yes "))
			return
		}
	}
	println("Node IP:   " + color.New(color.BgRed, color.Bold).Sprint(" No "))
}
