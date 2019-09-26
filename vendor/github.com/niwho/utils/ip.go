package utils

import "net"

func LookupIP() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ip := ipnet.IP.To4()
			if ip != nil {
				return ip.String()
			}
		}
	}
	return ""
}

func Hashing() int64 {
	ip := LookupIP()
	if ip == "" {
		return 0
	}
	return 0
}
