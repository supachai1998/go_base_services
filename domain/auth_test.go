package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestTokenExpires_Expired(t *testing.T) {
	type fields struct {
		BaseModel BaseModel
		Token     string
		ExpireAt  time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Expired",
			fields: fields{
				ExpireAt: time.Now().Add(-1 * time.Second),
			},
			want: true,
		},
		{
			name: "Not Expired",
			fields: fields{
				ExpireAt: time.Now().Add(1 * time.Second),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := TokenExpires{
				BaseModel: tt.fields.BaseModel,
				Token:     tt.fields.Token,
				ExpireAt:  tt.fields.ExpireAt,
			}
			if got := rt.Expired(); got != tt.want {
				t.Errorf("TokenExpires.Expired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {
	type args struct {
		claims jwt.Claims
		secret string
	}
	tests := []struct {
		name          string
		args          args
		wantStartWith string
		wantErr       bool
	}{
		{
			name: "Ok",
			args: args{
				claims: func() jwt.Claims {
					now := time.Now()
					return &VerifyTokenClaims{
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
							IssuedAt:  jwt.NewNumericDate(now),
							NotBefore: jwt.NewNumericDate(now),
						},
						TokenType: StaffVerifyTokenType,
						Token:     "token",
					}
				}(),
				secret: "secret",
			},
			wantErr:       false,
			wantStartWith: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateAccessToken(tt.args.claims, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.HasPrefix(got, tt.wantStartWith) {
				t.Errorf("GenerateAccessToken() = %v, want %v", got, tt.wantStartWith)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	correctToken := func() string {
		claims := VerifyTokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
			TokenType: StaffVerifyTokenType,
			Token:     "token",
		}
		token, _ := GenerateAccessToken(claims, "secret")
		return token
	}
	correctClaims := func() jwt.Claims {
		now := time.Now()
		return &VerifyTokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
			},
			TokenType: StaffVerifyTokenType,
			Token:     "token",
		}
	}()
	type args struct {
		tokenStr string
		claims   jwt.Claims
		secret   string
	}
	tests := []struct {
		name    string
		args    args
		want    *jwt.Token
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				tokenStr: correctToken(),
				claims:   correctClaims,
				secret:   "secret",
			},
			wantErr: false,
			want: &jwt.Token{
				Raw: correctToken(),
			},
		},
		{
			name: "Expired",
			args: args{
				tokenStr: func() string {
					claims := VerifyTokenClaims{
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
							IssuedAt:  jwt.NewNumericDate(time.Now()),
							NotBefore: jwt.NewNumericDate(time.Now()),
						},
						TokenType: StaffVerifyTokenType,
						Token:     "token",
					}
					token, _ := GenerateAccessToken(claims, "secret")
					return token
				}(),
				claims: correctClaims,
			},
			wantErr: true,
		},
		{
			name: "Method not match",
			args: args{
				tokenStr: func() string {
					// Use a different claims type
					claims := jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						NotBefore: jwt.NewNumericDate(time.Now()),
					}
					token := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
					tokenStr, _ := token.SignedString([]byte("secret"))
					return tokenStr
				}(),
				claims: correctClaims,
			},
			wantErr: true,
		},
		{
			name: "NotBefore Exceeded",
			args: args{
				tokenStr: func() string {
					claims := VerifyTokenClaims{
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
							IssuedAt:  jwt.NewNumericDate(time.Now()),
							NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
						},
						TokenType: StaffVerifyTokenType,
						Token:     "token",
					}

					token, _ := GenerateAccessToken(claims, "secret")
					return token
				}(),
				claims: correctClaims,
			},
			wantErr: true,
		},
		{
			name: "Invalid Signature",
			args: args{
				tokenStr: func() string {
					claims := VerifyTokenClaims{
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
							IssuedAt:  jwt.NewNumericDate(time.Now()),
							NotBefore: jwt.NewNumericDate(time.Now()),
						},
						TokenType: StaffVerifyTokenType,
						Token:     "token",
					}
					token, _ := GenerateAccessToken(claims, "wrongsecret")
					return token
				}(),
				claims: func() jwt.Claims {
					now := time.Now()
					return &VerifyTokenClaims{
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
							IssuedAt:  jwt.NewNumericDate(now),
							NotBefore: jwt.NewNumericDate(now),
						},
						TokenType: StaffVerifyTokenType,
						Token:     "token",
					}
				}(),
				secret: "secret",
			},
			wantErr: true,
		},
		{
			name: "Invalid Token",
			args: args{
				tokenStr: "abc",
				claims:   &jwt.MapClaims{},
				secret:   "secret",
			},
			wantErr: true,
		},
		{
			name: "Invalid Token Method",
			args: args{
				tokenStr: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJSUzI1NmluT1RBIiwibmFtZSI6IkpvaG4gRG9lIn0.ICV6gy7CDKPHMGJxV80nDZ7Vxe0ciqyzXD_Hr4mTDrdTyi6fNleYAyhEZq2J29HSI5bhWnJyOBzg2bssBUKMYlC2Sr8WFUas5MAKIr2Uh_tZHDsrCxggQuaHpF4aGCFZ1Qc0rrDXvKLuk1Kzrfw1bQbqH6xTmg2kWQuSGuTlbTbDhyhRfu1WDs-Ju9XnZV-FBRgHJDdTARq1b4kuONgBP430wJmJ6s9yl3POkHIdgV-Bwlo6aZluophoo5XWPEHQIpCCgDm3-kTN_uIZMOHs2KRdb6Px-VN19A5BYDXlUBFOo-GvkCBZCgmGGTlHF_cWlDnoA9XTWWcIYNyUI4PXNw",
				claims:   correctClaims,
				secret:   "secret",
			},
			wantErr: true,
		},
		{
			name: "Invalid Expiration Time",
			args: args{
				tokenStr: func() string {
					// Use jwt.MapClaims and set an invalid expiration time
					claims := jwt.MapClaims{
						"exp": "invalid",
					}
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					tokenStr, _ := token.SignedString([]byte("secret"))
					return tokenStr
				}(),
				claims: correctClaims,
				secret: "secret",
			},
			wantErr: true,
		},
		{
			name: "Expired Expiration Time",
			args: args{
				tokenStr: func() string {
					// Use jwt.MapClaims and set an invalid expiration time
					claims := jwt.MapClaims{
						"exp": int64(time.Now().Add(-time.Second).Unix()),
					}
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					tokenStr, _ := token.SignedString([]byte("secret"))
					return tokenStr
				}(),
				claims: correctClaims,
				secret: "secret",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseToken(tt.args.tokenStr, tt.args.claims, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got.Raw != tt.want.Raw {
				t.Errorf("ParseToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
