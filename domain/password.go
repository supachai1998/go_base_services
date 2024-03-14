package domain

import (
	"go_base/hash"
	"go_base/logger"
)

type Password string

// CompareBcrypt
func (p Password) CompareBcrypt(password string) bool {
	p = Password(hash.Trim(string(p)))
	if p == "" || password == "" {
		logger.L().Warn("password is nil")
		return false
	}
	return hash.CompareBcrypt(string(p), password)
}

func (p Password) Hash() Password {
	p = Password(hash.Trim(string(p)))
	if p == "" {
		logger.L().Warn("password is nil")
		return ""
	}
	hash, _ := hash.HashBcrypt(string(p))
	return Password(hash)
}

func (p Password) String() string {
	return string(p)
}

func (p Password) IsHashed() bool {
	return hash.IsHashed(string(p))
}
