package utils

import (
	"CrestedIbis/src/global"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtToken struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenToken 生成token信息
func (JwtToken) GenToken(username string) (string, error) {
	claims := JwtToken{
		username,
		[]string{"admin"},
		jwt.RegisteredClaims{
			Issuer: "CrestedIbis",
			// 过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.Conf.Jwt.ExpireTime) * time.Second)),
			// 生效时间
			NotBefore: jwt.NewNumericDate(time.Now()),
			// 签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(global.Conf.Jwt.Key))
}

func (JwtToken) GenTempAccessToken(username string, roles []string, expireTime uint) (string, error) {
	claims := JwtToken{
		username,
		roles,
		jwt.RegisteredClaims{
			Issuer: "CrestedIbis",
			// 过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireTime) * time.Second)),
			// 生效时间
			NotBefore: jwt.NewNumericDate(time.Now()),
			// 签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(global.Conf.Jwt.Key))
}

func (JwtToken) ParseToken(tokenString string) (*JwtToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Conf.Jwt.Key), nil
	})

	if claims, ok := token.Claims.(*JwtToken); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
