package utils

import (
	"time"

	"pluginmarket-server/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int `json:"userId"`
	Power  int `json:"power"`
	jwt.RegisteredClaims
}

// GenerateToken 生成普通用户 JWT
func GenerateToken(userID, power int) (string, error) {
	claims := Claims{
		UserID: userID,
		Power:  power,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.C.JWT.Expire) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.C.JWT.Secret))
}

// GenerateShortToken 生成短期 JWT（用于邮箱验证、密码重置，10分钟有效）
func GenerateShortToken(userID int, email, tokenType string, version int64) (string, error) {
	claims := jwt.MapClaims{
		"userId":  userID,
		"email":   email,
		"type":    tokenType,
		"version": version,
		"exp":     time.Now().Add(10 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.C.JWT.Secret))
}

// ParseToken 解析普通 JWT
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.C.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

// ParseShortToken 解析短期 JWT
func ParseShortToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.C.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
