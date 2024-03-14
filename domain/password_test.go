package domain

import (
	"strings"
	"testing"
)

func TestPassword_CompareBcrypt(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		p    Password
		args args
		want bool
	}{
		{
			name: "pass hash",
			args: args{
				password: "test",
			},
			p:    Password("test"),
			want: false,
		},
		{
			name: "pass is nil",
			args: args{
				password: "test",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.CompareBcrypt(tt.args.password); got != tt.want {
				t.Errorf("Password.CompareBcrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassword_Hash(t *testing.T) {
	tests := []struct {
		name string
		p    Password
		want string
	}{
		{
			name: "p is empty",
			p:    Password(""),
			want: "",
		},
		{
			name: "p ",
			p:    Password("a"),
			want: "$2a$10$",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Hash(); !strings.Contains(got.String(), string(tt.want)) {
				t.Errorf("Password.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassword_String(t *testing.T) {
	tests := []struct {
		name string
		p    Password
		want string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("Password.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassword_IsHashed(t *testing.T) {
	tests := []struct {
		name string
		p    Password
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.IsHashed(); got != tt.want {
				t.Errorf("Password.IsHashed() = %v, want %v", got, tt.want)
			}
		})
	}
}
