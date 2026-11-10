package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess = "access"
)

type TokenManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type Claims struct {
	jwt.RegisteredClaims

	UserID uint32 `json:"user_id"`
}

func NewTokenManager(
	privateKey *rsa.PrivateKey,
	publicKey *rsa.PublicKey,
	issuer string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *TokenManager {
	return &TokenManager{
		privateKey: privateKey,
		publicKey:  publicKey,
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (m *TokenManager) GenerateAccessToken(user *models.User) (models.AccessToken, error) {
	now := time.Now()

	claims := Claims{
		// #nosec G115 -- user.ID is a postgres INTEGER and fits in uint32
		UserID: uint32(user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   TokenTypeAccess,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenStr, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", err
	}

	return models.AccessToken(tokenStr), nil
}

func (m *TokenManager) ValidateAccessToken(tokenStr string) (*domain.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodRS256.Alg() {
				return nil, domain.ErrTokenInvalid
			}
			return m.publicKey, nil
		},
	)

	if err != nil || !token.Valid {
		return nil, domain.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, domain.ErrTokenInvalid
	}

	if claims.Issuer != m.issuer {
		return nil, domain.ErrTokenInvalid
	}

	if claims.Subject != TokenTypeAccess {
		return nil, domain.ErrTokenInvalid
	}

	return &domain.AccessTokenClaims{
		UserID:    claims.UserID,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

func (m *TokenManager) GenerateRefreshToken() (*models.RefreshToken, error) {
	bytes := make([]byte, 64)

	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	raw := hex.EncodeToString(bytes)

	return &models.RefreshToken{
		Value:     raw,
		Hash:      m.HashToken(raw),
		ExpiresAt: time.Now().Add(m.refreshTTL),
	}, nil
}

func (m *TokenManager) HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
