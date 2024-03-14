package domain

// Sensitive Data
type SensitiveString string

func (s SensitiveString) String() string {
	return string(s)
}
