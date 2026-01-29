package main

// Local Mode

import (
	_ "embed"
	"fmt"
	"net/netip"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

//go:embed ip2region_v4.xdb
var ip2RegionV4XDBData []byte

//go:embed ip2region_v6.xdb
var ip2RegionV6XDBData []byte

func getIpLocationByIp2Region(ipStr string) string {
	addr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return fmt.Sprintf("failed to parse ip: %s", err)
	}

	var version *xdb.Version
	var cBuff []byte
	if addr.Is4() {
		version = xdb.IPv4
		cBuff = ip2RegionV4XDBData
	} else if addr.Is6() {
		version = xdb.IPv6
		cBuff = ip2RegionV6XDBData
	} else {
		return fmt.Sprintf("invalid ip: %s", ipStr)
	}

	searcher, err := xdb.NewWithBuffer(version, cBuff)
	if err != nil {
		return fmt.Sprintf("failed to create searcher: %s", err)
	}
	defer searcher.Close()

	region, err := searcher.SearchByStr(ipStr)
	if err != nil {
		return fmt.Sprintf("failed to search: %s", err)
	}

	return fmt.Sprintf("%s|ip2region|local", region)
}
