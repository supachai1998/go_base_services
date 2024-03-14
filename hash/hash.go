package hash

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	CharsetV1 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	CharsetV2 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*_+"
	CharsetV3 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+{}[]:;\"'<>?,./|\\"

	bcryptRegex = "^\\$2(?:a|b|x|y)\\$(?:[4-9]|[12][0-9]|3[01])\\$(?:[a-zA-Z\\d\\./]{22})(?:[a-zA-Z\\d\\./]{31})$"
)

// bcrypt
func HashBcrypt(str string, costOpt ...int) (string, error) {
	bCost := bcrypt.DefaultCost
	if len(costOpt) > 0 {
		bCost = costOpt[0]
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(str), bCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CompareBcrypt(hashedPassword, str string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(str)) == nil
}

func Trim(password string) string {
	return strings.ReplaceAll(password, " ", "")
}

func GenerateRandomString(length int, charsetOpt ...string) string {
	if length <= 0 {
		return ""
	}
	charset := CharsetV2
	if len(charsetOpt) > 0 {
		charset = charsetOpt[0]
	}
	if length > 128 {
		length = 128
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[RandomInt(len(charset))]
	}
	return string(b)
}

func RandomInt(max int) int {
	if max <= 0 {
		return 0
	}
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))
	return ran.Intn(max)
}

func GenerateToken() string {
	return GenerateRandomString(32)
}

func IsHashed(password string) bool {
	if ok, _ := regexp.MatchString(bcryptRegex, password); ok {
		return true
	}
	return false
}
