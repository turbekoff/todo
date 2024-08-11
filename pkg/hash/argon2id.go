package hash

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Hasher interface {
	Hash(password string) (string, error)
	Compare(password string, hash string) error
}

var (
	ErrInvalid      = errors.New("invalid hash format")
	ErrMismatched   = errors.New("hash doesn't match")
	ErrIncompatible = errors.New("incompatible version")
)

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type Argon2idHasher struct {
	memory     uint32
	time       uint32
	threads    uint8
	keyLength  uint32
	saltLength uint32
	pepper     string
}

func NewArgon2idHasher(pepper string) Hasher {
	return &Argon2idHasher{
		memory:     46 * 1024,
		time:       1,
		threads:    1,
		keyLength:  32,
		saltLength: 16,
		pepper:     pepper,
	}
}

func (h *Argon2idHasher) Hash(password string) (string, error) {
	salt, err := generateRandomBytes(h.saltLength)
	if err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(password+h.pepper), salt, h.time, h.memory, h.threads, h.keyLength)
	content := base64.RawStdEncoding.EncodeToString(append(key, salt...))
	return fmt.Sprintf("v=%d$%s", argon2.Version, content), nil
}

func (h *Argon2idHasher) Compare(password string, hash string) error {
	parts := strings.Split(hash, "$")
	if len(parts) != 2 {
		return ErrInvalid
	}

	var version int
	if _, err := fmt.Sscanf(parts[0], "v=%d", &version); err != nil {
		return ErrIncompatible
	}

	content, err := base64.RawStdEncoding.Strict().DecodeString(parts[1])
	if err != nil {
		return err
	}

	key := []byte(content)[:h.keyLength]
	salt := []byte(content)[h.keyLength:]
	passwordKey := argon2.IDKey([]byte(password+h.pepper), salt, h.time, h.memory, h.threads, h.keyLength)

	if string(passwordKey) != string(key) {
		return ErrMismatched
	}
	return nil
}
