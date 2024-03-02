package url

import "strings"

func SameHost(host string, expected string) bool {
	return strings.EqualFold(host, expected)
}