package domain

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestBaseModel_IsZeroID(t *testing.T) {
	type fields struct {
		ID        uuid.UUID
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *gorm.DeletedAt
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Is Zero ID",
			fields: fields{
				ID: uuid.Nil,
			},
			want: true,
		},
		{
			name: "Isn't Zero ID",
			fields: fields{
				ID: uuid.New(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &BaseModel{
				ID:        tt.fields.ID,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
				DeletedAt: tt.fields.DeletedAt,
			}
			if got := m.IsZeroID(); got != tt.want {
				t.Errorf("BaseModel.IsZeroID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseModel_IsDeleted(t *testing.T) {
	type fields struct {
		ID        uuid.UUID
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *gorm.DeletedAt
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Is Deleted",
			fields: fields{
				DeletedAt: &gorm.DeletedAt{
					Valid: true,
					Time:  time.Now(),
				},
			},
			want: true,
		},
		{
			name: "Is Deleted",
			fields: fields{
				DeletedAt: &gorm.DeletedAt{
					Valid: false,
					Time:  time.Now(),
				},
			},
			want: true,
		},
		{
			name: "Is not Deleted",
			fields: fields{
				DeletedAt: &gorm.DeletedAt{
					Valid: false,
					Time:  time.Time{},
				},
			},
			want: false,
		},
		{
			name: "Is not Deleted",
			fields: fields{
				DeletedAt: &gorm.DeletedAt{
					Valid: true,
					Time:  time.Time{},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &BaseModel{
				ID:        tt.fields.ID,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
				DeletedAt: tt.fields.DeletedAt,
			}
			if got := m.IsDeleted(); got != tt.want {
				t.Errorf("BaseModel.IsDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseModel_BeforeCreate(t *testing.T) {
	type fields struct {
		ID        uuid.UUID
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *gorm.DeletedAt
	}
	type args struct {
		tx *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Is Zero ID",
			fields: fields{},
			args: args{
				tx: &gorm.DB{},
			},
			wantErr: false,
		},
		{
			name: "Isn't Zero ID",
			fields: fields{
				ID: uuid.New(),
			},
			args: args{
				tx: &gorm.DB{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &BaseModel{
				ID:        tt.fields.ID,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
				DeletedAt: tt.fields.DeletedAt,
			}
			if err := m.BeforeCreate(tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("BaseModel.BeforeCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBaseModel_BeforeUpdate(t *testing.T) {
	type fields struct {
		ID        uuid.UUID
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *gorm.DeletedAt
	}
	type args struct {
		tx *gorm.DB
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		wantUpdatedAt *time.Time
	}{
		{
			name:   "Is Zero ID",
			fields: fields{},
			args: args{
				tx: &gorm.DB{},
			},
			wantErr: false,
		},
		{
			name: "Isn't Zero ID",
			fields: fields{
				ID: uuid.New(),
			},
			args: args{
				tx: &gorm.DB{},
			},
			wantErr:       false,
			wantUpdatedAt: func() *time.Time { t := time.Now(); return &t }(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &BaseModel{
				ID:        tt.fields.ID,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
				DeletedAt: tt.fields.DeletedAt,
			}
			if err := m.BeforeUpdate(tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("BaseModel.BeforeUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantUpdatedAt != nil {
				if m.UpdatedAt.Sub(*tt.wantUpdatedAt) > time.Second {
					t.Errorf("BaseModel.BeforeUpdate() UpdatedAt = %v, wantUpdatedAt %v, diff %v", m.UpdatedAt, tt.wantUpdatedAt, m.UpdatedAt.Sub(*tt.wantUpdatedAt))
				}
			}
		})
	}
}

func TestConvertAnyIntoBaseModel(t *testing.T) {
	type args struct {
		modelAny any
	}
	tests := []struct {
		name string
		args args
		want BaseModel
	}{
		{
			name: "Is BaseModel",
			args: args{
				modelAny: BaseModel{},
			},
			want: BaseModel{},
		},
		{
			name: "Is BaseModel Ptr",
			args: args{
				modelAny: &BaseModel{},
			},
			want: BaseModel{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertAnyIntoBaseModel(tt.args.modelAny); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertAnyIntoBaseModel() = %v, want %v", got, tt.want)
			}
		})
	}
}
