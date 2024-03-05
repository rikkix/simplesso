package crypto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	// Secret is the secret key for the JWT.
	secret []byte
}

// NewAuth creates a new Auth instance.
func NewAuth(secret []byte) *Auth {
	return &Auth{
		secret: secret,
	}
}

// GenerateToken generates a new JWT token.
func (a *Auth) GenerateToken(sub string, aud string, exp time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   sub,
		Audience:  jwt.ClaimStrings{aud},
		ExpiresAt: jwt.NewNumericDate(exp),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secret)
}

// ValidateToken validates a JWT token.
func (a *Auth) ValidateToken(token string, aud string) (string, time.Time, error) {
	tk, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return a.secret, nil
	}, jwt.WithExpirationRequired(), jwt.WithAudience(aud))
	if err != nil {
		return "", time.Time{}, err
	}

	sub, err := tk.Claims.GetSubject()
	if err != nil {
		return "", time.Time{}, err
	}

	exp, err := tk.Claims.GetExpirationTime()

	return sub, exp.Time, err
}
