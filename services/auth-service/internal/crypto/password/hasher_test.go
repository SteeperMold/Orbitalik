package password

import (
	"errors"
	"strings"
	"testing"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
)

func TestHasher_Hash(t *testing.T) {
	hasher := NewHasher()

	// #nosec G101 -- definitely not a real password :)
	password := "WHY IS THERE CODE??? MAKE A *** .EXE FILE AND GIVE IT TO ME"

	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	if hash == "" {
		t.Fatal("Hash() returned empty hash")
	}

	if !strings.HasPrefix(hash, "$argon2id$v=19$") {
		t.Fatalf("Hash() produced unexpected format: %q", hash)
	}

	if err := hasher.Compare(hash, password); err != nil {
		t.Fatalf("Compare() failed for correct password: %v", err)
	}
}

func TestHasher_Hash_UsesDifferentSalt(t *testing.T) {
	hasher := NewHasher()

	password := "same password"

	hash1, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("first Hash() error = %v", err)
	}

	hash2, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("second Hash() error = %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("Hash() returned identical hashes for the same password")
	}
}

func TestHasher_Compare(t *testing.T) {
	hasher := NewHasher()

	hash, err := hasher.Hash("correct password")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	tests := []struct {
		name     string
		hash     string
		password string
		wantErr  error
	}{
		{
			name:     "correct password",
			hash:     hash,
			password: "correct password",
			wantErr:  nil,
		},
		{
			name:     "wrong password",
			hash:     hash,
			password: "wrong password",
			wantErr:  domain.ErrInvalidCredentials,
		},
		{
			name:     "empty password",
			hash:     hash,
			password: "",
			wantErr:  domain.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hasher.Compare(tt.hash, tt.password)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"Compare() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestHasher_Compare_InvalidHash(t *testing.T) {
	hasher := NewHasher()

	tests := []struct {
		name string
		hash string
	}{
		{
			name: "empty hash",
			hash: "",
		},
		{
			name: "random string",
			hash: "not a hash",
		},
		{
			name: "wrong number of parts",
			hash: "$argon2id$v=19",
		},
		{
			name: "wrong algorithm",
			hash: "$bcrypt$v=19$m=65536,t=3,p=4$YWJj$YWJj",
		},
		{
			name: "wrong version",
			hash: "$argon2id$v=18$m=65536,t=3,p=4$YWJj$YWJj",
		},
		{
			name: "invalid parameters",
			hash: "$argon2id$v=19$invalid$YWJj$YWJj",
		},
		{
			name: "malformed argon2 hash",
			hash: "$argon2id$v=19$invalid",
		},
		{
			name: "invalid salt encoding",
			hash: "$argon2id$v=19$m=65536,t=3,p=4$!!!$abc",
		},
		{
			name: "invalid hash encoding",
			hash: "$argon2id$v=19$m=65536,t=3,p=4$YWJj$!!!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hasher.Compare(tt.hash, "password")

			if !errors.Is(err, domain.ErrInvalidCredentials) {
				t.Fatalf(
					"Compare() error = %v, want %v",
					err,
					domain.ErrInvalidCredentials,
				)
			}
		})
	}
}

func TestHasher_Compare_InvalidParameters(t *testing.T) {
	hasher := NewHasher()

	hash, err := hasher.Hash("password")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	tests := []struct {
		name string
		hash string
	}{
		{
			name: "memory is zero",
			hash: strings.Replace(hash, "m=65536,", "m=0,", 1),
		},
		{
			name: "memory is too large",
			hash: strings.Replace(hash, "m=65536,", "m=65537,", 1),
		},
		{
			name: "time is zero",
			hash: strings.Replace(hash, "t=3,", "t=0,", 1),
		},
		{
			name: "time is too large",
			hash: strings.Replace(hash, "t=3,", "t=4,", 1),
		},
		{
			name: "threads is zero",
			hash: strings.Replace(hash, "p=4", "p=0", 1),
		},
		{
			name: "threads is too large",
			hash: strings.Replace(hash, "p=4", "p=5", 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hasher.Compare(tt.hash, "password")

			if !errors.Is(err, domain.ErrInvalidCredentials) {
				t.Fatalf(
					"Compare() error = %v, want %v",
					err,
					domain.ErrInvalidCredentials,
				)
			}
		})
	}
}
