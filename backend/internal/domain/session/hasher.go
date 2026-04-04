package session

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	Hash(s string) (string, error)
	Compare(hash, s string) error
}

type BcryptHasher struct{ Cost int }

func NewBcryptHasher(cost int) PasswordHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{Cost: cost}
}

func (a *BcryptHasher) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), a.Cost)
	if err != nil {
		return "", nil
	}
	return string(b), err
}

func (a *BcryptHasher) Compare(hash, s string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
}
