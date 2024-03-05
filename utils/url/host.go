package url

import (
	"strings"

	"github.com/valyala/fasthttp"
)

func SameHost(host string, expected string) bool {
	return strings.EqualFold(host, expected)
}

func ExtractHost(url string) string {
	uri := fasthttp.URI{}
	uri.Parse(nil, []byte(url))
	return string(uri.Host())
}