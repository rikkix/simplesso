package session

import (
	"github.com/valyala/fasthttp"
)

type Session struct {
	Method string
	URI *fasthttp.URI
	IP string
	UserAgent string
	Authorized bool
	Sub string
}

func (s *Session) Query(key string) []byte {
	return s.URI.QueryArgs().Peek(key)
}