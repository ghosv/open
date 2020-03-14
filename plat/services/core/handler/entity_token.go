package handler

import (
	"time"

	"github.com/ghosv/open/meta"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ghosv/open/plat/conf"
	pb "github.com/ghosv/open/plat/services/core/proto"
)

// Token of JWT
type Token struct {
	jwt.StandardClaims // TODO: switch to micro's jwt ?
	Payload            *pb.TokenPayload
}

// NewToken by userinfo
func NewToken() *Token {
	return &Token{
		StandardClaims: jwt.StandardClaims{},
	}
}

// Sign JWT
func (c *Token) Sign(d time.Duration) (string, error) {
	iat := time.Now()
	expireTime := iat.Add(d).Unix()
	c.Issuer = meta.SystemName
	c.IssuedAt = iat.Unix()
	c.ExpiresAt = expireTime
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return jwtToken.SignedString(conf.ENV.KeyJWT())
}

// FreshToken if will expire
func FreshToken(t *Token, freshD time.Duration, termD time.Duration) (string, error) {
	if time.Now().Add(freshD).After(time.Unix(t.ExpiresAt, 0)) {
		return t.Sign(termD)
	}
	return "", nil
}

// ValidToken for UserInfo
func ValidToken(tokenStr string) (t *Token, e error) {
	if tokenStr == "" {
		return nil, meta.ErrInvalidJWT
	}
	token, err := jwt.ParseWithClaims(tokenStr, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return conf.ENV.KeyJWT(), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, meta.ErrInvalidJWT
	}

	if t, ok := token.Claims.(*Token); ok {
		return t, nil
	}
	return nil, meta.ErrInvalidJWT
}
