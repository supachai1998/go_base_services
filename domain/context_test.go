package domain

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestUserID(t *testing.T) {
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Isn't UserID",
			args: args{
				ctx: func() echo.Context {
					e := echo.New()
					c := e.NewContext(nil, nil)
					c.Set(string(UserIDKey), "123")
					return c
				}(),
			},
			want: "",
		},
		{
			name: "Is UserID",
			args: args{
				ctx: func() echo.Context {
					e := echo.New()
					c := e.NewContext(nil, nil)
					c.Set(string(UserIDKey), "05ec8667-b1f0-4771-b74d-64af4678c45f")
					return c
				}(),
			},
			want: "05ec8667-b1f0-4771-b74d-64af4678c45f",
		},
		{
			name: "No UserID",
			args: args{
				ctx: echo.New().NewContext(nil, nil),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserID(tt.args.ctx); got != tt.want {
				t.Errorf("UserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStaffFromContext(t *testing.T) {
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name string
		args args
		want *Staff
	}{
		{
			name: "Is Staff",
			args: args{
				ctx: func() echo.Context {
					e := echo.New()
					c := e.NewContext(nil, nil)
					c.Set(string(StaffKey), &Staff{})
					return c
				}(),
			},
			want: &Staff{},
		},
		{
			name: "No Staff",
			args: args{
				ctx: echo.New().NewContext(nil, nil),
			},
			want: nil,
		},
		{
			name: "Staff is not Staff",
			args: args{
				ctx: func() echo.Context {
					e := echo.New()
					c := e.NewContext(nil, nil)
					c.Set(string(StaffKey), "123")
					return c
				}(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StaffFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StaffFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetActionFromContext(t *testing.T) {
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Is Get Action",
			args: args{
				ctx: func() echo.Context {
					e := echo.New()
					c := e.NewContext(nil, nil)
					c.SetPath("/api/v1/staffs")
					c.SetRequest(&http.Request{
						Method: http.MethodGet,
					})
					return c
				}(),
			},
			want: "GET|staffs",
		},
		{
			name: "Is Post Action",
			args: args{
				ctx: func() echo.Context {
					e := echo.New()
					c := e.NewContext(nil, nil)
					c.SetPath("/api/v1/staffs")
					c.SetRequest(&http.Request{
						Method: http.MethodPost,
					})
					return c
				}(),
			},
			want: "POST|staffs",
		},
		{
			name: "Is Put Action",
			args: args{
				ctx: func() echo.Context {
					e := echo.New()
					c := e.NewContext(nil, nil)
					c.SetPath("/api/v1/staffs")
					c.SetRequest(&http.Request{
						Method: http.MethodPut,
					})
					return c
				}(),
			},
			want: "PUT|staffs",
		},
		{
			name: "Is Delete Action",
			args: args{
				ctx: func() echo.Context {
					e := echo.New()
					c := e.NewContext(nil, nil)
					c.SetPath("/api/v1/staffs")
					c.SetRequest(&http.Request{
						Method: http.MethodDelete,
					})
					return c
				}(),
			},
			want: "DELETE|staffs",
		},
		{
			name: "Is Patch Action",
			args: args{
				ctx: func() echo.Context {
					e := echo.New()
					c := e.NewContext(nil, nil)
					c.SetPath("/api/v1/staffs")
					c.SetRequest(&http.Request{
						Method: http.MethodPatch,
					})
					return c
				}(),
			},
			want: "PATCH|staffs",
		},
		{
			name: "No Action",
			args: args{
				ctx: echo.New().NewContext(nil, nil),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetActionFromContext(tt.args.ctx); got != tt.want {
				t.Errorf("GetActionFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrLogGlsGo(t *testing.T) {
	type args struct {
		ctx echo.Context
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Is Error",
			args: args{
				ctx: echo.New().NewContext(&http.Request{}, nil),
				err: echo.ErrNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("ErrLogGlsGo() did panic server down is failed")
					}
				}()
				ErrLogGlsGo(tt.args.ctx, tt.args.err)
			}()
		})
	}
}

func TestGetUUIDFromParam(t *testing.T) {
	type args struct {
		ctx echo.Context
		key string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 uuid.UUID
	}{
		{
			name: "Is UUID",
			args: args{
				ctx: func() echo.Context {
					e := echo.New()
					c := e.NewContext(nil, nil)
					c.SetParamNames("id")
					c.SetParamValues("05ec8667-b1f0-4771-b74d-64af4678c45f")
					return c
				}(),
				key: "id",
			},
			want:  "05ec8667-b1f0-4771-b74d-64af4678c45f",
			want1: uuid.MustParse("05ec8667-b1f0-4771-b74d-64af4678c45f"),
		},
		{
			name: "No UUID",
			args: args{
				ctx: echo.New().NewContext(nil, nil),
				key: "id",
			},
			want:  "",
			want1: uuid.Nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetUUIDFromParam(tt.args.ctx, tt.args.key)
			if got != tt.want {
				t.Errorf("GetUUIDFromParam() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetUUIDFromParam() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetUUID(t *testing.T) {
	type args struct {
		idStr string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 uuid.UUID
	}{
		{
			name: "Is UUID",
			args: args{
				idStr: "05ec8667-b1f0-4771-b74d-64af4678c45f",
			},
			want:  "05ec8667-b1f0-4771-b74d-64af4678c45f",
			want1: uuid.MustParse("05ec8667-b1f0-4771-b74d-64af4678c45f"),
		},
		{
			name: "No UUID",
			args: args{
				idStr: "123",
			},
			want:  "",
			want1: uuid.Nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetUUID(tt.args.idStr)
			if got != tt.want {
				t.Errorf("GetUUID() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetUUID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
