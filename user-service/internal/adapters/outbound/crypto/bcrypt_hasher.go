package crypto

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

type Config struct {
	HashCost int
}

func (bc *BcryptHasher) Hash(plain string) (string, error) {
	start := time.Now()
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bc.cost)
	log.Printf("bcrypt.Hash took %v", time.Since(start))
	return string(b), err
}

func (bc *BcryptHasher) Compare(hash string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

func NewBcryptHasher(cfg Config) *BcryptHasher {
	return &BcryptHasher{cost: cfg.HashCost}
}
