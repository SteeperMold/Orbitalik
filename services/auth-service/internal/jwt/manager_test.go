package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/jwt"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

const testIssuer = "orbitalik-auth"

func newTestManager(t *testing.T) *jwt.TokenManager {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	return jwt.NewTokenManager(
		privateKey,
		&privateKey.PublicKey,
		testIssuer,
		time.Hour,
		7*24*time.Hour,
	)
}

func newTestUser() *models.User {
	return &models.User{
		ID:       42,
		Email:    "test@example.com",
		Username: "testuser",
	}
}

func TestTokenManager_GenerateAccessToken(t *testing.T) {
	manager := newTestManager(t)
	user := newTestUser()

	token, err := manager.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	if token == "" {
		t.Fatal("GenerateAccessToken() returned empty token")
	}

	parts := strings.Split(string(token), ".")
	if len(parts) != 3 {
		t.Fatalf("JWT has %d parts, want 3", len(parts))
	}

	claims, err := manager.ValidateAccessToken(string(token))
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	// #nosec G115 -- user.ID is 42 which fits in uint32
	if claims.UserID != uint32(user.ID) {
		t.Errorf(
			"UserID = %d, want %d",
			claims.UserID,
			user.ID,
		)
	}

	if claims.ExpiresAt.IsZero() {
		t.Error("ExpiresAt is zero")
	}

	if !claims.ExpiresAt.After(time.Now()) {
		t.Error("ExpiresAt is not in the future")
	}
}

func TestTokenManager_GenerateAccessToken_ContainsExpectedClaims(t *testing.T) {
	manager := newTestManager(t)
	user := newTestUser()

	token, err := manager.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	parts := strings.Split(string(token), ".")
	if len(parts) != 3 {
		t.Fatalf("invalid JWT format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("decode payload: %v", err)
	}

	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		t.Fatalf("unmarshal claims: %v", err)
	}

	if got := claims["iss"]; got != testIssuer {
		t.Errorf("iss = %v, want %q", got, testIssuer)
	}

	if got := claims["sub"]; got != jwt.TokenTypeAccess {
		t.Errorf("sub = %v, want %q", got, jwt.TokenTypeAccess)
	}

	if got := claims["user_id"]; got != float64(user.ID) {
		t.Errorf("user_id = %v, want %d", got, user.ID)
	}

	if _, ok := claims["iat"]; !ok {
		t.Error("iat claim is missing")
	}

	if _, ok := claims["exp"]; !ok {
		t.Error("exp claim is missing")
	}
}

func TestTokenManager_GenerateAccessToken_Expiration(t *testing.T) {
	accessTTL := 2 * time.Hour

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	manager := jwt.NewTokenManager(
		privateKey,
		&privateKey.PublicKey,
		testIssuer,
		accessTTL,
		time.Hour,
	)

	before := time.Now().Add(accessTTL)

	token, err := manager.GenerateAccessToken(newTestUser())
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	claims, err := manager.ValidateAccessToken(string(token))
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	after := time.Now().Add(accessTTL)

	if claims.ExpiresAt.Before(before.Add(-time.Second)) ||
		claims.ExpiresAt.After(after.Add(time.Second)) {
		t.Errorf(
			"ExpiresAt = %v, expected approximately now + %v",
			claims.ExpiresAt,
			accessTTL,
		)
	}
}

func TestTokenManager_ValidateAccessToken(t *testing.T) {
	manager := newTestManager(t)

	token, err := manager.GenerateAccessToken(newTestUser())
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	claims, err := manager.ValidateAccessToken(string(token))
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("UserID = %d, want 42", claims.UserID)
	}
}

func TestTokenManager_ValidateAccessToken_InvalidTokens(t *testing.T) {
	manager := newTestManager(t)

	// #nosec G101 -- fake tokens
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty",
			token: "",
		},
		{
			name:  "random string",
			token: "not-a-jwt",
		},
		{
			name:  "malformed jwt",
			token: "abc.def",
		},
		{
			name:  "invalid signature",
			token: "eyJhbGciOiJSUzI1NiJ9.eyJ1c2VyX2lkIjoxfQ.invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := manager.ValidateAccessToken(tt.token)

			if !errors.Is(err, domain.ErrTokenInvalid) {
				t.Fatalf(
					"ValidateAccessToken() error = %v, want %v",
					err,
					domain.ErrTokenInvalid,
				)
			}
		})
	}
}

func TestTokenManager_ValidateAccessToken_WrongSigningKey(t *testing.T) {
	manager := newTestManager(t)

	token, err := manager.GenerateAccessToken(newTestUser())
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	otherKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	otherManager := jwt.NewTokenManager(
		otherKey,
		&otherKey.PublicKey,
		testIssuer,
		time.Hour,
		time.Hour,
	)

	_, err = otherManager.ValidateAccessToken(string(token))
	if !errors.Is(err, domain.ErrTokenInvalid) {
		t.Fatalf(
			"ValidateAccessToken() error = %v, want %v",
			err,
			domain.ErrTokenInvalid,
		)
	}

	_ = manager
}

func TestTokenManager_ValidateAccessToken_Expired(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	manager := jwt.NewTokenManager(
		privateKey,
		&privateKey.PublicKey,
		testIssuer,
		-1*time.Second,
		time.Hour,
	)

	token, err := manager.GenerateAccessToken(newTestUser())
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	_, err = manager.ValidateAccessToken(string(token))
	if !errors.Is(err, domain.ErrTokenInvalid) {
		t.Fatalf(
			"ValidateAccessToken() error = %v, want %v",
			err,
			domain.ErrTokenInvalid,
		)
	}
}

func TestTokenManager_ValidateAccessToken_WrongAlgorithm(t *testing.T) {
	manager := newTestManager(t)
	user := newTestUser()

	claims := jwtlib.MapClaims{
		"iss":     testIssuer,
		"sub":     jwt.TokenTypeAccess,
		"user_id": user.ID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwtlib.NewWithClaims(
		jwtlib.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString([]byte("wrong-secret"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	_, err = manager.ValidateAccessToken(tokenString)
	if !errors.Is(err, domain.ErrTokenInvalid) {
		t.Fatalf(
			"ValidateAccessToken() error = %v, want %v",
			err,
			domain.ErrTokenInvalid,
		)
	}
}

func TestTokenManager_ValidateAccessToken_WrongIssuer(t *testing.T) {
	manager := newTestManager(t)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	claims := jwtlib.MapClaims{
		"iss":     "attacker",
		"sub":     jwt.TokenTypeAccess,
		"user_id": 42,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwtlib.NewWithClaims(
		jwtlib.SigningMethodRS256,
		claims,
	)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	_, err = manager.ValidateAccessToken(tokenString)

	if !errors.Is(err, domain.ErrTokenInvalid) {
		t.Errorf(
			"ValidateAccessToken() error = %v, want %v",
			err,
			domain.ErrTokenInvalid,
		)
	}
}

func TestTokenManager_ValidateAccessToken_WrongSubject(t *testing.T) {
	manager := newTestManager(t)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	claims := jwtlib.MapClaims{
		"iss":     testIssuer,
		"sub":     "refresh",
		"user_id": 42,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwtlib.NewWithClaims(
		jwtlib.SigningMethodRS256,
		claims,
	)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	_, err = manager.ValidateAccessToken(tokenString)

	if !errors.Is(err, domain.ErrTokenInvalid) {
		t.Errorf(
			"ValidateAccessToken() error = %v, want %v",
			err,
			domain.ErrTokenInvalid,
		)
	}
}

func TestTokenManager_GenerateRefreshToken(t *testing.T) {
	manager := newTestManager(t)

	token, err := manager.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	if token == nil {
		t.Fatal("GenerateRefreshToken() returned nil")
	}

	if token.Value == "" {
		t.Fatal("refresh token Value is empty")
	}

	if token.Hash == "" {
		t.Fatal("refresh token Hash is empty")
	}

	if token.ExpiresAt.IsZero() {
		t.Fatal("refresh token ExpiresAt is zero")
	}

	if !token.ExpiresAt.After(time.Now()) {
		t.Fatal("refresh token is already expired")
	}

	if token.Hash != manager.HashToken(token.Value) {
		t.Error("refresh token Hash doesn't match Value")
	}
}

func TestTokenManager_GenerateRefreshToken_Unique(t *testing.T) {
	manager := newTestManager(t)

	token1, err := manager.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("first GenerateRefreshToken() error = %v", err)
	}

	token2, err := manager.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("second GenerateRefreshToken() error = %v", err)
	}

	if token1.Value == token2.Value {
		t.Fatal("two refresh tokens have identical values")
	}

	if token1.Hash == token2.Hash {
		t.Fatal("two refresh tokens have identical hashes")
	}
}

func TestTokenManager_GenerateRefreshToken_Length(t *testing.T) {
	manager := newTestManager(t)

	token, err := manager.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	if len(token.Value) != 128 {
		t.Errorf(
			"refresh token length = %d, want 128",
			len(token.Value),
		)
	}

	if len(token.Hash) != 64 {
		t.Errorf(
			"refresh token hash length = %d, want 64",
			len(token.Hash),
		)
	}
}

func TestTokenManager_GenerateRefreshToken_ExpiresAt(t *testing.T) {
	refreshTTL := 24 * time.Hour

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	manager := jwt.NewTokenManager(
		privateKey,
		&privateKey.PublicKey,
		testIssuer,
		time.Hour,
		refreshTTL,
	)

	before := time.Now().Add(refreshTTL)

	token, err := manager.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	after := time.Now().Add(refreshTTL)

	if token.ExpiresAt.Before(before.Add(-time.Second)) ||
		token.ExpiresAt.After(after.Add(time.Second)) {
		t.Errorf(
			"ExpiresAt = %v, expected approximately now + %v",
			token.ExpiresAt,
			refreshTTL,
		)
	}
}

func TestTokenManager_HashToken(t *testing.T) {
	manager := newTestManager(t)

	// #nosec G101 -- fake tokens
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "normal token",
			token: "abc123",
		},
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "unicode",
			token: "токен🔐",
		},
		{
			name:  "long token",
			token: strings.Repeat("a", 10000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.HashToken(tt.token)

			if got == "" {
				t.Fatal("HashToken() returned empty hash")
			}

			if len(got) != 64 {
				t.Errorf(
					"HashToken() length = %d, want 64",
					len(got),
				)
			}

			gotAgain := manager.HashToken(tt.token)

			if got != gotAgain {
				t.Error("HashToken() is not deterministic")
			}
		})
	}
}

func TestTokenManager_HashToken_DifferentInputs(t *testing.T) {
	manager := newTestManager(t)

	hash1 := manager.HashToken("token-1")
	hash2 := manager.HashToken("token-2")

	if hash1 == hash2 {
		t.Fatal("different tokens produced identical hashes")
	}
}
