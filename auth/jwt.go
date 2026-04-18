package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Issuer string        `yaml:"issuer"`
	Secret string        `yaml:"secret"`
	Expire time.Duration `yaml:"expire"`
}

type Claims struct {
	UserID    int64  `json:"user_id"`
	SessionID int64  `json:"session_id"`
	UserType  string `json:"user_type"`
	Username  string `json:"username"`
	jwt.RegisteredClaims
}

type IssuePayload struct {
	UserID    int64
	SessionID int64
	UserType  string
	Username  string
	Expire    time.Duration
}

type Manager struct {
	cfg JWTConfig
}

func NewManager(cfg JWTConfig) *Manager {
	return &Manager{cfg: cfg}
}

func (m *Manager) Issue(payload IssuePayload) (string, *Claims, error) {
	expire := payload.Expire
	if expire <= 0 {
		expire = m.cfg.Expire
	}
	if expire <= 0 {
		expire = 24 * time.Hour
	}

	now := time.Now()
	claims := &Claims{
		UserID:    payload.UserID,
		SessionID: payload.SessionID,
		UserType:  payload.UserType,
		Username:  payload.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.cfg.Issuer,
			Subject:   fmt.Sprintf("%d", payload.UserID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(m.cfg.Secret))
	if err != nil {
		return "", nil, fmt.Errorf("sign jwt: %w", err)
	}
	return signed, claims, nil
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {
	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(m.cfg.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse jwt: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid jwt token")
	}
	return claims, nil
}
