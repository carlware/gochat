package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Profile struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	PasswordHashed string `json:"-"`
}

func NewProfile(name, avatar, password string) *Profile {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return &Profile{
		ID:             uuid.New().String(),
		Name:           name,
		Avatar:         avatar,
		PasswordHashed: string(hashedPassword),
	}
}

func (p *Profile) PasswordMatch(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHashed), []byte(password))
	if err == nil {
		return true
	}
	return false
}
