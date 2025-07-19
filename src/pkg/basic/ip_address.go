package basic

import (
	"fmt"
	"net"
	"sync/atomic"
)

var (
	// 本机IP地址
	ipAddress atomic.Value

	// 私有网络地址段
	// 10.0.0.0 - 10.255.255.255（掩码范围需在16 - 28之间）
	// 172.16.0.0 - 172.31.255.255（掩码范围需在16 - 28之间）
	// 192.168.0.0 - 192.168.255.255 （掩码范围需在16 - 28之间）
	_, cidr10, _  = net.ParseCIDR("10.0.0.0/8")
	_, cidr172, _ = net.ParseCIDR("172.0.0.0/8")
	_, cidr192, _ = net.ParseCIDR("192.168.0.0/16")
)

func GetIPAddress() (string, error) {
	if v := ipAddress.Load(); v != nil {
		return v.(string), nil
	} else if addresses, err := net.InterfaceAddrs(); err != nil {
		return "", err
	} else {
		for _, addr := range addresses {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					ipStr := ipNet.IP.String()
					ipAddress.Store(ipStr)
					return ipStr, nil
				}
			}
		}
		return "", fmt.Errorf("not found")
	}
}

func IsPrivateIP(ip string) bool {
	parseIP := net.ParseIP(ip).To4()

	switch {
	case parseIP == nil:
		return false
	case cidr10.Contains(parseIP), cidr192.Contains(parseIP):
		return true
	case cidr172.Contains(parseIP) && parseIP[1] >= 16 && parseIP[1] <= 31:
		return true
	default:
		return false
	}
}

func CheckCidrIntersect(cidrs []string) (bool, error) {
	nets := []*net.IPNet{}
	for _, cidr := range cidrs {
		_, item, err := net.ParseCIDR(cidr)
		if err != nil {
			return false, err
		}
		for i := range nets {
			if nets[i].Contains(item.IP) || item.Contains(nets[i].IP) {
				return true, nil
			}
		}
		nets = append(nets, item)
	}
	return false, nil
}
