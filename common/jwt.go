package common

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Jwts struct {
	JwtKey     string
	ExpireHour int
}

type Claims struct {
	UserId int
	jwt.StandardClaims
}

func NewJwtw(j JwtConf) *Jwts {
	return &Jwts{
		JwtKey:     j.key,
		ExpireHour: j.expireHour,
	}
}

// 颁发token
func (j Jwts) GenerateToken(userId int) (string, error) {
	expireTime := time.Now().Add(time.Hour * time.Duration(j.ExpireHour))
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
	tokenStr, err := token.SignedString([]byte(j.JwtKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// 解析token
func (j Jwts) ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(j.JwtKey), nil
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
func (j Jwts) RefreshToken(tokenStr string) (string, error) {
	claims, err := j.ParseToken(tokenStr)
	if err != nil {
		return "", err
	}

	token, err := j.GenerateToken(claims.UserId)
	if err != nil {
		return "", err
	}

	return token, nil
}
