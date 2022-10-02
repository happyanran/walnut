package common

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserId uint32
	jwt.StandardClaims
}

// 颁发token
func (s ServiceContext) GenerateToken(userId uint32) (string, error) {
	expireTime := time.Now().Add(time.Hour * time.Duration(s.Cfg.JwtConf.expireHour))
	claims := Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(), // 签名时间
			Issuer:    "freedb.com",      // 签名颁发者
			Subject:   "user token",      // 签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.Cfg.JwtConf.key))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// 解析token
func (s ServiceContext) ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(s.Cfg.JwtConf.key), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// 刷新 Token
func (s ServiceContext) RefreshToken(tokenStr string) (string, error) {
	claims, err := s.ParseToken(tokenStr)
	if err != nil {
		return "", err
	}

	token, err := s.GenerateToken(claims.UserId)
	if err != nil {
		return "", err
	}

	return token, nil
}
