package validator

import (
	"fmt"
	"net"
	"strings"
)

// ValidateIPAddress validates an IP address with optional CIDR notation
func ValidateIPAddress(ip string) error {
	// If CIDR notation is used, separate IP address and network mask parts
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

	// For regular IP addresses
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address format")
	}
	return nil
}

// ValidateDNSAddress validates a DNS server address
func ValidateDNSAddress(dns string) error {
	// DNS server address must be an IPv4 or IPv6 address
	return ValidateIPAddress(dns)
}
