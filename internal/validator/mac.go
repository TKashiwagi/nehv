package validator

import (
	"errors"
	"regexp"
)

var macRegex = regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)

// ValidateMACAddress checks if the given string is a valid MAC address (colon-separated)
func ValidateMACAddress(mac string) error {
	if !macRegex.MatchString(mac) {
		return errors.New("invalid MAC address format (expected XX:XX:XX:XX:XX:XX)")
	}
	return nil
}
