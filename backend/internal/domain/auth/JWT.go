package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrAuthSecretRequired    = errors.New("jwt secret is required")
	ErrAuthIssuerRequired    = errors.New("jwt issuer is required")
	ErrAuthAccessTTLRequired = errors.New("jwt accessTTL is required")
)

type JWTConfig struct {
	Secret    string
	Issuer    string
	AccessTTL time.Duration
}

type JWTManager struct {
	Secret     []byte
	Issuer     string
	AccesssTTL time.Duration
}

type AccessClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTManager(cfg JWTConfig) (*JWTManager, error) {
	if cfg.Secret == "" {
		return nil, ErrAuthSecretRequired
	}

	if cfg.Issuer == "" {
		return nil, ErrAuthIssuerRequired
	}

	ttl := cfg.AccessTTL
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	return &JWTManager{
		Secret:     []byte(cfg.Secret),
		Issuer:     cfg.Issuer,
		AccesssTTL: cfg.AccessTTL,
	}, nil
}

func (a *JWTManager) NewAccessToken(userID, role string) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(a.AccesssTTL)

	claims := AccessClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    a.Issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString(a.Secret)
	return s, exp, err
}
