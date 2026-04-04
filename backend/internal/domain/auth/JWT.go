package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrAuthSecretRequired    = errors.New("jwt secret is required")
	ErrAuthSecretEmpty       = errors.New("jwt secret can't be empty")
	ErrAuthIssuerRequired    = errors.New("jwt issuer is required")
	ErrAuthClock             = errors.New("jwt clock skew must be >= 0")
	ErrAuthAccessTTLRequired = errors.New("jwt accessTTL is required")
	ErrAuthUserRequired      = errors.New("userID is required")
	ErrMissingToken          = errors.New("missing token")
	ErrInvalidSignAlg        = errors.New("invalid signing algorithm")
	ErrTokenExpired          = errors.New("token has expired")
	ErrTokenNotValid         = errors.New("token is not valid")
)

type JWTConfig struct {
	Secret    string
	Issuer    string
	Audience  string
	AccessTTL time.Duration
	ClockSkew time.Duration
}

type JWTManager struct {
	Secret     []byte
	Issuer     string
	Audience   string
	AccesssTTL time.Duration
	ClockSkew  time.Duration
}

type AccessClaims struct {
	Role   string   `json:"role,omitempty"`
	Scopes []string `json:"scopes,omitempty"`

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

	skew := cfg.ClockSkew
	if skew < 0 {
		return nil, ErrAuthClock
	}

	if skew == 0 {
		skew = 30 * time.Second
	}

	return &JWTManager{
		Secret:     []byte(cfg.Secret),
		Issuer:     cfg.Issuer,
		Audience:   cfg.Audience,
		AccesssTTL: ttl,
		ClockSkew:  skew,
	}, nil
}

func (j *JWTManager) NewAccessToken(userID, role string, scopes []string) (string, time.Time, error) {
	if len(j.Secret) == 0 {
		return "", time.Time{}, ErrAuthSecretEmpty
	}

	if userID == "" {
		return "", time.Time{}, ErrAuthUserRequired
	}

	now := time.Now()
	exp := now.Add(j.AccesssTTL)

	claims := AccessClaims{
		Role:   role,
		Scopes: scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			Subject:   userID,
			Audience:  audIfSet(j.Audience),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(-j.ClockSkew)),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString(j.Secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return s, exp, err
}

// Parses a JWT string,  verifies signature/alg, validates exp/nbf/iat, and returns strongly typed claims
func (j *JWTManager) ParseAndValidate(tokenSTR string) (*AccessClaims, error) {
	tokenStr := strings.TrimSpace(tokenSTR)
	if tokenStr == "" {
		return nil, ErrMissingToken
	}

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(j.Issuer),
		jwt.WithLeeway(j.ClockSkew),
	)

	claims := &AccessClaims{}
	token, err := parser.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidSignAlg
		}
		return j.Secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValid
		}
		return nil, err
	}

	if !token.Valid {
		return nil, ErrTokenNotValid
	}

	if strings.TrimSpace(claims.Subject) == "" {
		return nil, ErrTokenNotValid
	}

	return claims, nil
}

// Extracts the token from the authHeader
func ExtractBearerToken(authHeader string) (string, error) {
	h := strings.TrimSpace(authHeader)
	if h == "" {
		return "", ErrMissingToken
	}

	const prefix = "Bearer "
	if len(h) < len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
		return "", ErrMissingToken
	}

	tiki := strings.TrimSpace(h[len(prefix):])
	if tiki == "" {
		return "", ErrMissingToken
	}

	return tiki, nil
}

///

func audIfSet(aud string) jwt.ClaimStrings {
	if strings.TrimSpace(aud) == "" {
		return nil
	}
	return jwt.ClaimStrings{aud}
}
