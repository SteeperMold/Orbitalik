package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"golang.org/x/crypto/argon2"
)

const (
	defaultTime    = 3
	defaultMemory  = 64 * 1024 // 64 MB
	defaultThreads = 4
	defaultLength  = 32
	defaultSaltLen = 16
)

type Hasher struct {
	time    uint32
	memory  uint32
	threads uint8
	length  uint32
	saltLen uint32
}

func NewHasher() *Hasher {
	return &Hasher{
		time:    defaultTime,
		memory:  defaultMemory,
		threads: defaultThreads,
		length:  defaultLength,
		saltLen: defaultSaltLen,
	}
}

func (h *Hasher) Hash(password string) (string, error) {
	salt := make([]byte, h.saltLen)

	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.time,
		h.memory,
		h.threads,
		h.length,
	)

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		h.memory,
		h.time,
		h.threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func (h *Hasher) Compare(passwordHash string, password string) error {
	parts := strings.Split(passwordHash, "$")

	if len(parts) != 6 {
		return domain.ErrInvalidCredentials
	}

	if parts[1] != "argon2id" {
		return domain.ErrInvalidCredentials
	}

	if parts[2] != "v=19" {
		return domain.ErrInvalidCredentials
	}

	var memory, timeCost, threads uint32

	_, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&memory,
		&timeCost,
		&threads,
	)
	if err != nil {
		return domain.ErrInvalidCredentials
	}

	if memory == 0 || memory > defaultMemory {
		return domain.ErrInvalidCredentials
	}

	if timeCost == 0 || timeCost > defaultTime {
		return domain.ErrInvalidCredentials
	}

	if threads == 0 || threads > defaultThreads {
		return domain.ErrInvalidCredentials
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return domain.ErrInvalidCredentials
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return domain.ErrInvalidCredentials
	}

	if len(expectedHash) != int(h.length) {
		return domain.ErrInvalidCredentials
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		timeCost,
		memory,
		uint8(threads),
		h.length,
	)

	if subtle.ConstantTimeCompare(expectedHash, actualHash) == 0 {
		return domain.ErrInvalidCredentials
	}

	return nil
}
