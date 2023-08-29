package main

import (
	"flag"
	"net"
	"strconv"
	"strings"
	"time"
)

var nodeCIDRsString = ""

func getNodeCIDRs(nodeCIDRsString string) (string, []*net.IPNet, error) {
	var nodeDate string
	nodeCIDRsSlice := strings.Split(nodeCIDRsString, ",")
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

func printNodeCIDRs(nodeCIDRsString string) {
	nodeCIDRsSlice := strings.Split(nodeCIDRsString, ",")
	println("Node List:")
	for i, nodeCIDR := range nodeCIDRsSlice {
		if i > 0 {
			println("  " + strconv.Itoa(i) + ". " + nodeCIDR)
		}
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
	nodeFile := flag.String("f", "", "ip address")
	printNodes := flag.Bool("p", false, "print node ip list")
	flag.Parse()
	return *ipStr, *nodeFile, *printNodes
}

func main() {
	nodeDate, nodeCIDRs, err := getNodeCIDRs(nodeCIDRsString)
	if err != nil {
		println("node.txt format error, the first line must be a date (format 'YYYY-mm-dd') and the following lines must be legal CIDRs format")
		return
	}

	println("Node IP Update At: " + nodeDate + ", A tool to check if an IP is a knownsec node ip")
	println()

	ipStr, _, printNodes := initFlag()
	if printNodes {
		printNodeCIDRs(nodeCIDRsString)
		return
	}
	if ipStr == "" {
		flag.Usage()
		return
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		println("ip address format error")
		return
	}

	for _, nodeCIDR := range nodeCIDRs {
		if nodeCIDR.Contains(ip) {
			println(ipStr + " IS a knownsec node ip")
			return
		}
	}
	println(ipStr + " IS NOT a knownsec node ip")
}
