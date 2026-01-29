package main

// Local Mode

import (
	"embed"
	"fmt"

	"github.com/lionsoul2014/ip2region/binding/golang/service"
)

//go:embed ip2region_v4.xdb
var ip2RegionV4XDB embed.FS

//go:embed ip2region_v6.xdb
var ip2RegionV6XDB embed.FS

func getIpLocationByIp2Region(ipStr string) string {
	v4Config, err := service.NewV4Config(service.BufferCache, "ip2region_v4.xdb", 20)
	if err != nil {
		return fmt.Sprintf("failed to create v4 config: %s", err)
	}

	// 2, 创建 v6 的配置：指定缓存策略和 v6 的 xdb 文件路径
	v6Config, err := service.NewV6Config(service.BufferCache, "ip2region_v6.xdb", 20)
	if err != nil {
		return fmt.Sprintf("failed to create v6 config: %s", err)
	}

	// 3，通过上述配置创建 Ip2Region 查询服务
	ip2region, err := service.NewIp2Region(v4Config, v6Config)
	if err != nil {
		return fmt.Sprintf("failed to create ip2region service: %s", err)
	}
	defer ip2region.Close()

	region, err := ip2region.SearchByStr(ipStr)
	if err != nil {
		return fmt.Sprintf("failed to search: %s", err)
	}

	return fmt.Sprintf("%s|ip2region|local", region)
}
