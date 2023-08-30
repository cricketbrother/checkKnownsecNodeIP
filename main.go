package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
)

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
	nodeFile := flag.String("f", "", "ip address")
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

	println("checkKnownsecNodeIP " + version + ", A tool to check if an IP is a knownsec node ip")
	println("Nodes Update At: " + nodeDate)
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

	for _, nodeCIDR := range nodeCIDRs {
		if nodeCIDR.Contains(ip) {
			println(color.New(color.BgGreen, color.Bold).Sprint("[√]") + " " + ipStr + " IS a knownsec node ip")
			return
		}
	}
	println(color.New(color.BgRed, color.Bold).Sprint("[X]") + " " + ipStr + " IS NOT a knownsec node ip")
}
