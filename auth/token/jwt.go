package token

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	"time"
)

type JWTTokenGen struct {
	Issuer     string
	nowFunc    func() time.Time
	privateKey *rsa.PrivateKey
}

func NewJWTTokenGen(issuer string, privateKey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		Issuer:     issuer,
		nowFunc:    time.Now,
		privateKey: privateKey,
	}
}

func (g *JWTTokenGen) GenerateToken(accountID string, expireTime time.Duration) (string, error) {
	nowSec := g.nowFunc().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    g.Issuer,
		IssuedAt:  nowSec,
		ExpiresAt: nowSec + int64(expireTime),
		Subject:   accountID,
	})
	return token.SignedString(g.privateKey)
}
