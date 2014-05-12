package auth

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
)

type (
	Authenticator interface {
		HashPassword(string) (string, error)
		Authenticate(hashed string, password string) bool
		GenerateToken() string
	}

	Auth struct {
		cost int
	}
)

func NewAuthenticator(cost int) Authenticator {
	auth := &Auth{
		cost: cost,
	}
	return auth
}

func (auth *Auth) HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), auth.cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (auth *Auth) Authenticate(hashed string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false
	}
	return true
}

func (auth *Auth) GenerateToken() string {
	return uuid.New()
}
