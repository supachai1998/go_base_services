package services_test

import (
	"fmt"
	"go_base/domain"
	"go_base/storage"

	"github.com/go-faker/faker/v4"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

var (
	failedPass  = "failed pass"
	successPass = "Passw0rd!"
)

func (uts *UnitTestSuite) TestStaffService_Unlock_IsVerified() {
	storage.Migration[domain.Role](uts.server.DB, "roles.json")
	err := storage.Seed[domain.StaffMigration](uts.server.DB, "staffs_unlock.json")
	if err != nil {
		uts.T().Fatal(err)
	}

	staffs, err := uts.service.Staff.Find(uts.ctx, domain.Pagination[domain.Staff]{
		PaginationSwagger: domain.PaginationSwagger{
			Page:     lo.ToPtr(1),
			PageSize: lo.ToPtr(10),
		},
	})
	if err != nil {
		uts.T().Fatal(err)
	}
	staff := staffs.Items[0]
	err = uts.server.Storages.Cache.DeleteStrikes(uts.ctx.Request().Context(), string(staff.Email))
	if err != nil {
		uts.T().Fatal(err)
	}
	for i := 0; i < uts.server.Cfg.AdminAuth.AccountLockoutMaxAttempts; i++ {
		_, err := uts.service.Staff.LoginWithEmailPassword(uts.ctx, domain.StaffLogin{
			Email:    staff.Email,
			Password: failedPass,
		})
		// should be invalid input
		if err == nil {
			uts.T().Fatal("should be invalid input")
		}
		uts.Equal("invalid_input", err.Error())

	}

	_, err = uts.service.Staff.LoginWithEmailPassword(uts.ctx, domain.StaffLogin{
		Email:    staff.Email,
		Password: failedPass,
	})
	// should be invalid input
	if err == nil {
		uts.T().Fatal("should be invalid input")
	}
	uts.Equal("forbidden", err.Error())

	if err = uts.service.Staff.Unlock(uts.ctx, domain.StaffUnlock{Email: string(staff.Email)}); err != nil {
		uts.T().Fatal(err)
	}
	strikes, err := uts.server.Storages.Cache.GetStrikes(uts.ctx.Request().Context(), string(staff.Email))
	if err != nil {
		uts.T().Fatal(err)
	}
	uts.Equal(0, strikes)

	success, err := uts.service.Staff.LoginWithEmailPassword(uts.ctx, domain.StaffLogin{
		Email:    staff.Email,
		Password: successPass,
	})

	if err != nil {
		uts.T().Fatal(err)
	}

	uts.NotEqual("", success.AccessToken)

}
func (uts *UnitTestSuite) TestStaffService_Unlock_IsNotVerified() {
	storage.Migration[domain.Role](uts.server.DB, "roles.json")
	err := storage.Seed[domain.StaffMigration](uts.server.DB, "staffs_unlock.json")
	if err != nil {
		uts.T().Fatal(err)
	}

	staffs, err := uts.service.Staff.Find(uts.ctx, domain.Pagination[domain.Staff]{
		PaginationSwagger: domain.PaginationSwagger{
			Page:     lo.ToPtr(1),
			PageSize: lo.ToPtr(10),
		},
	})
	if err != nil {
		uts.T().Fatal(err)
	}
	staff := staffs.Items[1]
	err = uts.server.Storages.Cache.DeleteStrikes(uts.ctx.Request().Context(), fmt.Sprintf(domain.StaffAuthCache, staff.Email))
	if err != nil {
		uts.T().Fatal(err)
	}
	for i := 0; i < uts.server.Cfg.AdminAuth.AccountLockoutMaxAttempts; i++ {
		_, err := uts.service.Staff.LoginWithEmailPassword(uts.ctx, domain.StaffLogin{
			Email:    staff.Email,
			Password: failedPass,
		})
		// should be invalid input
		if err == nil {
			uts.T().Fatal("should be invalid input")
		}
		uts.Equal("invalid_input", err.Error())

	}

	_, err = uts.service.Staff.LoginWithEmailPassword(uts.ctx, domain.StaffLogin{
		Email:    staff.Email,
		Password: failedPass,
	})
	// should be invalid input
	if err == nil {
		uts.T().Fatal("should be invalid input")
	}
	uts.Equal("forbidden", err.Error())
	err = uts.service.Staff.Unlock(uts.ctx, domain.StaffUnlock{
		Email: string(staff.Email),
	})
	if err != nil {
		uts.T().Fatal(err)
	}
	strikes, err := uts.server.Storages.Cache.GetStrikes(uts.ctx.Request().Context(), string(staff.Email))
	if err != nil {
		uts.T().Fatal(err)
	}
	uts.Equal(0, strikes)

	_, err = uts.service.Staff.LoginWithEmailPassword(uts.ctx, domain.StaffLogin{
		Email:    staff.Email,
		Password: "Passw0rd!",
	})

	if err == nil {
		uts.T().Fatal("should be invalid input")
	}
	uts.Equal("invalid_input", err.Error())

}
func (uts *UnitTestSuite) TestStaffService_Unlock_InvalidEmail() {
	storage.Migration[domain.Role](uts.server.DB, "roles.json")
	err := storage.Seed[domain.StaffMigration](uts.server.DB, "staffs_unlock.json")
	if err != nil {
		uts.T().Fatal(err)
	}

	err = uts.service.Staff.Unlock(uts.ctx, domain.StaffUnlock{
		Email: "invalid@email.com",
	})
	if err != nil {
		uts.T().Fatal(err)
	}

}

func (uts *UnitTestSuite) TestStaffService_Login_InvalidEmail() {
	storage.Migration[domain.Role](uts.server.DB, "roles.json")
	err := storage.Seed[domain.StaffMigration](uts.server.DB, "staffs_unlock.json")
	if err != nil {
		uts.T().Fatal(err)
	}
	err = uts.server.Storages.Cache.DeleteStrikes(uts.ctx.Request().Context(), fmt.Sprintf(domain.StaffAuthCache, "invalid email"))
	if err != nil {
		uts.T().Fatal(err)
	}
	_, err = uts.service.Staff.LoginWithEmailPassword(uts.ctx, domain.StaffLogin{
		Email:    "invalid email",
		Password: failedPass,
	})
	// should be invalid input
	if err == nil {
		uts.T().Fatal("should be invalid input")
	}
	uts.Equal("invalid_input", err.Error())

}

func (uts *UnitTestSuite) TestStaffService_Find() {
	storage.Migration[domain.Role](uts.server.DB, "roles.json")
	storage.Mock[domain.StaffMock](20, uts.server.DB, func() domain.StaffMock {
		return domain.StaffMock{
			Email:      domain.SensitiveString(faker.Email()),
			FirstName:  faker.FirstName(),
			LastName:   faker.LastName(),
			Password:   domain.Password("mockPassW0rd").Hash().String(),
			IsVerified: true,
		}
	})

	// Find staff default page 1, page size 10
	staffs, err := uts.service.Staff.Find(uts.ctx, domain.Pagination[domain.Staff]{
		PaginationSwagger: domain.PaginationSwagger{
			Page:     lo.ToPtr(1),
			PageSize: lo.ToPtr(10),
		},
	})
	if err != nil {
		uts.T().Fatal(err)
	}
	uts.Equal(10, len(staffs.Items))

	// Find staffs page 2, page size 10
	staffs, err = uts.service.Staff.Find(uts.ctx, domain.Pagination[domain.Staff]{
		PaginationSwagger: domain.PaginationSwagger{
			Page:     lo.ToPtr(2),
			PageSize: lo.ToPtr(10),
		},
	})
	if err != nil {
		uts.T().Fatal(err)
	}
	uts.Equal(10, len(staffs.Items))

	storage.Drop[domain.Staff](uts.server.DB)
	storage.Mock[domain.StaffMock](50, uts.server.DB, func() domain.StaffMock {
		return domain.StaffMock{
			Email:      domain.SensitiveString(faker.Email()),
			FirstName:  faker.FirstName(),
			LastName:   faker.LastName(),
			Password:   domain.Password("mockPassW0rd").Hash().String(),
			IsVerified: true,
		}
	})

	staffs, err = uts.service.Staff.Find(uts.ctx, domain.Pagination[domain.Staff]{
		PaginationSwagger: domain.PaginationSwagger{
			PageSize: lo.ToPtr(50),
		},
	})
	if err != nil {
		uts.T().Fatal(err)
	}
	uts.Equal(50, len(staffs.Items))
}

func (uts *UnitTestSuite) TestStaffService_Get() {
	storage.Migration[domain.Role](uts.server.DB, "roles.json")
	storage.Mock[domain.StaffMock](100, uts.server.DB, func() domain.StaffMock {
		return domain.StaffMock{
			Email:      domain.SensitiveString(faker.Email()),
			FirstName:  faker.FirstName(),
			LastName:   faker.LastName(),
			Password:   domain.Password("mockPassW0rd").Hash().String(),
			IsVerified: true,
		}
	})
	staffs, err := uts.service.Staff.Find(uts.ctx, domain.Pagination[domain.Staff]{
		PaginationSwagger: domain.PaginationSwagger{
			Page:     lo.ToPtr(1),
			PageSize: lo.ToPtr(10),
		},
	})
	if err != nil {
		uts.T().Fatal(err)
	}
	uts.Equal(10, len(staffs.Items))

	// Get staffs by id
	staffsGet, err := uts.service.Staff.Get(uts.ctx, staffs.Items[0].ID.String())
	if err != nil {
		uts.T().Fatal(err)
	}
	uts.Equal(staffs.Items[0].ID, staffsGet.ID)

	storage.Drop[domain.Staff](uts.server.DB)
	storage.Mock[domain.StaffMock](50, uts.server.DB, func() domain.StaffMock {
		return domain.StaffMock{
			Email:      domain.SensitiveString(faker.Email()),
			FirstName:  faker.FirstName(),
			LastName:   faker.LastName(),
			Password:   domain.Password("mockPassW0rd").Hash().String(),
			IsVerified: true,
		}
	})
	staffs, err = uts.service.Staff.Find(uts.ctx, domain.Pagination[domain.Staff]{
		PaginationSwagger: domain.PaginationSwagger{
			Page:     lo.ToPtr(1),
			PageSize: lo.ToPtr(10),
		},
	})
	if err != nil {
		uts.T().Fatal(err)
	}
	// Get staffs by id
	staffGet, err := uts.service.Staff.Get(uts.ctx, staffs.Items[0].ID.String())
	if err != nil {
		uts.T().Fatal(err)
	}
	uts.Equal(staffs.Items[0].ID, staffGet.ID)
}

func (uts *UnitTestSuite) TestStaffService_Create() {
	storage.Migration[domain.Role](uts.server.DB, "roles.json")
	err := storage.Drop[domain.Staff](uts.server.DB)
	if err != nil {
		uts.T().Fatal(err)
	}

	// Create staff
	staffCreate, err := uts.service.Staff.Create(uts.ctx, domain.StaffCreate{
		Email:     "add1@email.com",
		FirstName: "FirstName",
		LastName:  "LastName",
	})
	if err != nil {
		uts.T().Fatal(err)
	}

	uts.Equal("FirstName", staffCreate.FirstName)
	uts.Equal("LastName", staffCreate.LastName)
}
func (uts *UnitTestSuite) TestStaffService_CreateDub() {
	storage.Migration[domain.Role](uts.server.DB, "roles.json")
	err := storage.Seed[domain.StaffMigration](uts.server.DB, "staffs.json")
	if err != nil {
		uts.T().Fatal(err)
	}

	// Create staff
	_, err = uts.service.Staff.Create(uts.ctx, domain.StaffCreate{
		Email:     "admin1@admin.com", // duplicate email
		FirstName: "FirstName",
		LastName:  "LastName",
	})
	if err == nil {
		uts.T().Fatal("should be duplicate email")
	}

	uts.Equal(gorm.ErrDuplicatedKey.Error(), err.Error())
}

func (uts *UnitTestSuite) TestStaffService_Delete() {
	storage.Migration[domain.Role](uts.server.DB, "roles.json")
	err := storage.Seed[domain.StaffMigration](uts.server.DB, "staffs.json")
	if err != nil {
		uts.T().Fatal(err)
	}
	staffs, err := uts.service.Staff.Find(uts.ctx, domain.Pagination[domain.Staff]{
		PaginationSwagger: domain.PaginationSwagger{
			Page:     lo.ToPtr(1),
			PageSize: lo.ToPtr(1),
		},
	})
	if err != nil {
		uts.T().Fatal(err)
	}

	uts.Equal(1, len(staffs.Items))
	staff := staffs.Items[0]
	// attach staff id to path
	uts.ctx.SetParamNames("id")
	uts.ctx.SetParamValues(staff.ID.String())
	// Delete staff
	err = uts.service.Staff.Delete(uts.ctx)
	if err != nil {
		uts.T().Fatal(err)
	}
	// find new staff after delete
	staffs, err = uts.service.Staff.Find(uts.ctx, domain.Pagination[domain.Staff]{
		PaginationSwagger: domain.PaginationSwagger{
			Page:     lo.ToPtr(1),
			PageSize: lo.ToPtr(1),
		},
	})
	if err != nil {
		uts.T().Fatal(err)
	}
	uts.Equal(1, len(staffs.Items))
}
