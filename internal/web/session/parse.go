package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rikkix/simplesso/internal/config"
	"github.com/rikkix/simplesso/utils/crypto"
	"github.com/rikkix/simplesso/utils/url"
	"github.com/valyala/fasthttp"
)

type SessionParser struct {
	SsoHost string
	cgp *CGISessionParser
	ssp *SSOSessionParser
}

// NewSessionParser creates a new SessionParser instance.
// TODO: Use the provided configuration to create a new SessionParser instance.
func NewSessionParser(conf *config.Server) *SessionParser {
	return &SessionParser{
		SsoHost: conf.SsoHost,
		cgp: &CGISessionParser{
			MethodHead: "X-Method",
			SchemeHead: "X-Scheme",
			URIHead: "X-URI",
			IPHead: "X-Real-IP",
			UAHead: "X-User-Agent",
			TokenCookieName: "auth_token",
			Auther: crypto.NewAuth(conf.GetServicesSecretBytes()),
		},
		ssp: &SSOSessionParser{
			IPHead: "X-Real-IP",
			TokenCookieName: "auth_token",
			Auther: crypto.NewAuth(conf.GetSsoSecretBytes()),
		},
	}
}

func (sp *SessionParser) Parse(ctx *fiber.Ctx) *Session {
	if url.SameHost(ctx.Hostname(), sp.SsoHost) {
		return sp.ssp.ParseSession(ctx)
	}
	return sp.cgp.ParseSession(ctx)
}

func (sp *SessionParser) SSOAuther() *crypto.Auth {
	return sp.ssp.Auther
}

func (sp *SessionParser) CGIAuther() *crypto.Auth {
	return sp.cgp.Auther
}

type CGISessionParser struct {
	MethodHead string
	SchemeHead string
	URIHead string
	IPHead string
	UAHead string
	TokenCookieName string
	Auther *crypto.Auth
}

// ParseSession parses a session from the fiber.Ctx.
func (sp *CGISessionParser) ParseSession(ctx *fiber.Ctx) *Session {
	u := fasthttp.URI{}
	u.Parse([]byte(ctx.Hostname()), []byte(ctx.Get(sp.URIHead)))
	u.SetScheme(ctx.Get(sp.SchemeHead))
	s := &Session{
		Method:     ctx.Get(sp.MethodHead),
		URI:        &u,
		// Host:       ctx.Hostname(),
		IP:         ctx.Get(sp.IPHead),
		UserAgent:  ctx.Get(sp.UAHead),
		Authorized: false,
	}
	sub, exp, err := sp.Auther.ValidateToken(ctx.Cookies(sp.TokenCookieName), string(ctx.Hostname()))
	if err == nil {
		s.Authorized = true
		s.Sub = sub
		s.Exp = exp
	}

	return s
}

type SSOSessionParser struct {
	IPHead string
	TokenCookieName string
	Auther *crypto.Auth
}	

// ParseSession parses a session from the fiber.Ctx.
func (sp *SSOSessionParser) ParseSession(ctx *fiber.Ctx) *Session {
	s := &Session{
		Method:     ctx.Method(),
		URI:        ctx.Context().URI(),
		IP:         ctx.Get(sp.IPHead),
		UserAgent:  ctx.Get(fiber.HeaderUserAgent),
		Authorized: false,
	}
	sub, exp, err := sp.Auther.ValidateToken(ctx.Cookies(sp.TokenCookieName), ctx.Hostname())
	if err == nil {
		s.Authorized = true
		s.Sub = sub
		s.Exp = exp
	}
	return s
}

