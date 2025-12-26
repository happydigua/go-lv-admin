package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"go-lv-vue-admin/internal/global"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(global.LV_CONFIG.JWT.SigningKey),
	}
}

// Custom Claims
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UserId   uint
	Username string
	RoleId   uint
}

func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	bf, _ := ParseDuration(global.LV_CONFIG.JWT.BufferTime)
	ep, _ := ParseDuration(global.LV_CONFIG.JWT.ExpiresTime)

	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bf / time.Second), // Buffer time for renewal
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000 * time.Second)), // Effective time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),                  // Expiration time
			Issuer:    global.LV_CONFIG.Zap.Prefix,                             // Issuer
		},
	}
	return claims
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}

// ParseToken 全局解析 Token 函数，供中间件使用
func ParseToken(tokenString string) (*CustomClaims, error) {
	j := NewJWT()
	return j.ParseToken(tokenString)
}

// Helper to parse duration string like "7d", "1d"
// Since time.ParseDuration doesn't support "d", we need a simple wrapper if we want to support it,
// OR just expect "168h" in config.
// For simplicity, let's assume config uses "h" or "m" for Go standard duration, OR implement a simple converter.
// The config said "7d". Standard time.ParseDuration doesn't support 'd'.
// Let's implement a simple parser for 'd'.
func ParseDuration(d string) (time.Duration, error) {
	d = strings.TrimSpace(d)
	dr, err := time.ParseDuration(d)
	if err == nil {
		return dr, nil
	}
	if strings.Contains(d, "d") {
		index := strings.Index(d, "d")
		hour, _ := strconv.Atoi(d[:index])
		dr = time.Hour * 24 * time.Duration(hour)
		return dr, nil
	}
	dv, err := strconv.ParseInt(d, 10, 64)
	return time.Duration(dv), err
}
