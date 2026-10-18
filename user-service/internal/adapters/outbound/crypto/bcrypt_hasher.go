package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

type Config struct {
	HashCost int
}

func (bc *BcryptHasher) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bc.cost)
	return string(b), err
}

func (bc *BcryptHasher) Compare(hash string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

func NewBcryptHasher(cfg Config) *BcryptHasher {
	return &BcryptHasher{cost: cfg.HashCost}
}
