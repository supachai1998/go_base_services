package hash

import (
	"testing"
)

func TestHashBcrypt(t *testing.T) {
	type args struct {
		password string
		cost     int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Test case 2",
			args: args{
				password: "abc123",
			},
			wantErr: false,
		},
		{
			name: "Test case 3",
			args: args{
				password: "",
			},
			wantErr: false,
		},
		{
			name: "Test case cost zero",
			args: args{
				password: "",
				cost:     0,
			},
			wantErr: false,
		},
		{
			name: "Test case cost negative",
			args: args{
				password: "",
				cost:     -1,
			},
			wantErr: false,
		},
		{
			name: "Test case bcrypt: password length exceeds 72 bytes",
			args: args{
				password: "12345678901234567890123456789012345678901234567890123456789012345678901234567890",
				cost:     -1,
			},
			wantErr: true,
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := HashBcrypt(tt.args.password, tt.args.cost)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashBcrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCompareBcrypt(t *testing.T) {
	type args struct {
		hashedPassword string
		password       string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1",
			args: args{
				password: "password123",
			},
			want: false,
		},
		{
			name: "Test case 2",
			args: args{
				password: "abc123",
			},
			want: false,
		},
		{
			name: "Test case 3",
			args: args{
				password: "",
			},
			want: false,
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareBcrypt(tt.args.hashedPassword, tt.args.password); got != tt.want {
				t.Errorf("CompareBcrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test case 1",
			args: args{
				password: "password  123",
			},
			want: "password123",
		},
		{
			name: "Test case 2",
			args: args{
				password: "a bc123",
			},
			want: "abc123",
		},
		{
			name: "Test case 3",
			args: args{
				password: "",
			},
			want: "",
		},
		{
			name: "Test case 4",
			args: args{
				password: "ฟ หก ฟกห ฟก",
			},
			want: "ฟหกฟกหฟก",
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trim(tt.args.password); got != tt.want {
				t.Errorf("Trim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRandomString(t *testing.T) {
	type args struct {
		length  int
		charset string
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
	}{
		{
			name: "Test case 1",
			args: args{
				length: 10,
			},
			wantLen: 10,
		},
		{
			name: "Test case 2",
			args: args{
				length: 6,
			},
			wantLen: 6,
		},
		{
			name: "Test case 3",
			args: args{
				length: -1,
			},
			wantLen: 0,
		},
		{
			name: "Test case 4",
			args: args{
				length: 99999,
			},
			wantLen: 128,
		},
		{
			name: "Test case 4",
			args: args{
				length:  99999,
				charset: CharsetV3,
			},
			wantLen: 128,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.charset != "" {
				if got := GenerateRandomString(tt.args.length, tt.args.charset); len(got) != tt.wantLen {
					t.Errorf("GenerateRandomString() = %v, wantLen %v", got, tt.wantLen)
				}
			} else {
				if got := GenerateRandomString(tt.args.length); len(got) != tt.wantLen {
					t.Errorf("GenerateRandomString() = %v, wantLen %v", got, tt.wantLen)
				}
			}

		})
	}
}

func TestRandomInt(t *testing.T) {
	type args struct {
		max int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1",
			args: args{
				max: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomInt(tt.args.max); got < 0 || got > tt.args.max {
				t.Errorf("RandomInt() = %v, want %v", got, tt.args.max)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		wantLen int
	}{
		{
			name:    "Test case 1",
			wantLen: 32,
		},
		{
			name:    "Test case 2",
			wantLen: 32,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateToken(); len(got) != tt.wantLen {
				t.Errorf("GenerateToken() = %v, wantLen %v", got, tt.wantLen)
			}
		})
	}
}

func TestIsHashed(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Test case paint text",
			args{
				password: "password123",
			},
			false,
		},
		{
			"Test case paint text",
			args{
				password: "asdlasdl",
			},
			false,
		},
		{
			"Test case paint text",
			args{
				password: "12/16/2079",
			},
			false,
		},
		{
			"Test case paint text",
			args{
				password: "2d88:7f0e:25b0:428a:010f:aad0:82dc:d0e2",
			},
			false,
		},
		{
			"Test case paint text",
			args{
				password: "unool@hakre.im",
			},
			false,
		},
		{
			"Test case paint text",
			args{
				password: "#b89aa3",
			},
			false,
		},
		{
			"Test case encrypted text 1",
			args{
				password: "$2a$10$M3gA9ldicWuDex3GAEtK9uXZMJHFJdGf17cEPXCnBV0wr7tICLfgK",
			},
			true,
		},
		{
			"Test case encrypted text 2",
			args{
				password: "$2a$10$Ic4NGMqntJ..pJKZtCquw.EYC.cbvEJlXlCVlpTt2jNX.E56WAn1q",
			},
			true,
		},
		{
			"Test case encrypted text 3",
			args{
				password: "$2a$10$WuetI0GIYyHCG8bhEBlbJeaFGOdo7QmFR3097vdtVfnakpWL5jYWa",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHashed(tt.args.password); got != tt.want {
				t.Errorf("IsHashed() = %v, want %v", got, tt.want)
			}
		})
	}
}
