package myutils

import (
	"net"
	"strconv"
	"strings"
)

// Convert uint to net.IP http://www.outofmemory.cn
func Inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// Convert net.IP to int64 ,  http://www.outofmemory.cn
func Inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")
	var sum int64
	if len(bits) >= 4 {
		b0, _ := strconv.Atoi(bits[0])
		b1, _ := strconv.Atoi(bits[1])
		b2, _ := strconv.Atoi(bits[2])
		b3, _ := strconv.Atoi(bits[3])
		sum += int64(b0) << 24
		sum += int64(b1) << 16
		sum += int64(b2) << 8
		sum += int64(b3)
	}
	return sum
}

// Convert uint to net.IP http://www.outofmemory.cn
func Inet4_ntoa(ipnr uint32) string {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

// Convert net.IP to int64 ,  http://www.outofmemory.cn
func Inet4_aton(ipnr string) uint32 {
	bits := strings.Split(ipnr, ".")
	if len(bits) == 4 {
		b0, _ := strconv.Atoi(bits[0])
		b1, _ := strconv.Atoi(bits[1])
		b2, _ := strconv.Atoi(bits[2])
		b3, _ := strconv.Atoi(bits[3])
		var sum uint32

		sum += uint32(b0) << 24
		sum += uint32(b1) << 16
		sum += uint32(b2) << 8
		sum += uint32(b3)
		return sum
	}
	return 0
}

func AddrToId(addr string) uint64 {
	ip, port, err := net.SplitHostPort(addr)
	if err == nil {
		ipu32 := Inet_aton(net.ParseIP(ip))
		port, _ := strconv.ParseInt(port, 10, 64)
		return uint64(ipu32) + uint64(port)<<32
	}
	return 0
}
