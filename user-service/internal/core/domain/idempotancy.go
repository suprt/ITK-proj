package domain

type IdempotencyRecord struct {
	PayloadHash string
	UserID      string
}
