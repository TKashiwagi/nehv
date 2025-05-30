package validator

import (
	"fmt"
	"net"
	"strings"
)

// ValidateIPAddress validates an IP address with optional CIDR notation
func ValidateIPAddress(ip string) error {
	// CIDR表記の場合は、IPアドレス部分とネットマスク部分を分離
	if strings.Contains(ip, "/") {
		ip, cidr, err := net.ParseCIDR(ip)
		if err != nil {
			return fmt.Errorf("invalid CIDR notation: %v", err)
		}
		if ip == nil {
			return fmt.Errorf("invalid IP address in CIDR notation")
		}
		if cidr == nil {
			return fmt.Errorf("invalid network mask in CIDR notation")
		}
		return nil
	}

	// 通常のIPアドレスの場合
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address format")
	}
	return nil
}

// ValidateDNSAddress validates a DNS server address
func ValidateDNSAddress(dns string) error {
	// DNSサーバーのアドレスはIPv4またはIPv6アドレスである必要がある
	return ValidateIPAddress(dns)
}
