package auth

type PasswordHasher interface {
	Hash(s string) (string, error)
	Compare(hash, s string) error
}
