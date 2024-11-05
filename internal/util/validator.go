package util

import (
	"net/url"
	"strconv"
	"strings"
)

func IsURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	parts := strings.Split(u.Host, ".")
	if !(len(parts) == 2 || len(parts) == 3) {
		return false
	}
	for _, part := range parts {
		if part == "" {
			return false
		}
	}
	lasts := strings.Split(parts[len(parts)-1], ":")
	tld := lasts[0]
	if len(tld) < 2 || len(tld) > 63 {
		return false
	}
	if len(lasts) > 2 {
		return false
	}
	if len(lasts) == 2 {
		port, err := strconv.Atoi(lasts[1])
		if err != nil {
			return false
		}
		return port < 65536
	}
	return true
}
