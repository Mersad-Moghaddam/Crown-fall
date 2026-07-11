package random

import (
	"crypto/hmac"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

const SeedSize = 32

type Source interface {
	Seed() ([]byte, error)
}

type CryptoSource struct{}

func (CryptoSource) Seed() ([]byte, error) {
	seed := make([]byte, SeedSize)
	_, err := cryptorand.Read(seed)
	return seed, err
}

type FixedSource struct{ Value []byte }

func (source FixedSource) Seed() ([]byte, error) {
	if len(source.Value) == 0 {
		return nil, errors.New("fixed seed cannot be empty")
	}
	return append([]byte(nil), source.Value...), nil
}

func Commitment(seed []byte) string {
	sum := sha256.Sum256(seed)
	return hex.EncodeToString(sum[:])
}

func Derive(seed []byte, domain string) []byte {
	mac := hmac.New(sha256.New, seed)
	mac.Write([]byte("crownfall/v1/" + domain))
	return mac.Sum(nil)
}
