package main

import (
	"context"
	"go_base/domain"
	helper "go_base/domain/helper"
	"go_base/server"
	"go_base/storage"
	"os"
	"strconv"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		panic("mockCount is required")
	}
	mockCountStr := args[0]
	mockCount, _ := strconv.Atoi(mockCountStr)
	if mockCount < 1 {
		panic("mockCount must be greater than 0")
	}
	app, err := server.CreateApp(context.Background())
	if err != nil {
		panic(err)
	}

	if err := storage.Migration[domain.Role](app.DB, "roles.json"); err != nil {
		panic(err)
	}
	if err := storage.Migration[domain.StaffMigration](app.DB, "staffs.json"); err != nil {
		panic(err)
	}
	if err := storage.Migration[domain.UserMigration](app.DB, "users.json"); err != nil {
		panic(err)
	}

	maxGoCount := 2
	ch := make(chan int, maxGoCount)
	go func() {
		if err := storage.Mock(1000, app.DB, func() domain.StaffMock {
			return domain.StaffMock{
				Email:     domain.SensitiveString(faker.Email()),
				FirstName: faker.FirstName(),
				LastName:  faker.LastName(),
				Password:  domain.Password("mockPassW0rd").Hash().String(),
				RoleID: lo.ToPtr(uuid.MustParse(helper.RandomOneOfLength([]string{
					"31b22992-cc69-4bcd-a10f-dd6abf6d0073",
					"852db843-8626-408c-ae3f-20256872441e",
					"02b78806-3e14-5570-a619-1f3b4a655df7",
					"becfe7a7-fda5-596c-90e6-ab861014e9d8",
				}))),
				IsVerified: true,
				Status:     domain.Status(helper.RandomOneOfLength([]string{"active", "inactive", "pending"})),
				Phone:      lo.ToPtr(helper.RandomPhoneTH()),
			}
		}); err != nil {
			panic(err)
		}
		ch <- 1
	}()

	go func() {
		if err := storage.Mock(mockCount, app.DB, func() domain.User {
			return domain.User{
				Email:          domain.SensitiveString(faker.Email()),
				Password:       domain.Password("mockPassW0rd"),
				IsVerified:     true,
				BudgetBuy:      lo.ToPtr(helper.RandomMinMaxFloat64(100000, 10000000)),
				BudgetSell:     lo.ToPtr(helper.RandomMinMaxFloat64(100000, 10000000)),
				BudgetPerMonth: lo.ToPtr(helper.RandomMinMaxFloat64(100000, 500000)),
				FirstName:      faker.FirstName(),
				LastName:       faker.LastName(),
				Phone:          lo.ToPtr(helper.RandomPhoneTH()),
				Source:         lo.ToPtr(helper.RandomOneOfLength([]string{"facebook", "google", "line", "email", "website"})),
				StaffID: lo.ToPtr(uuid.MustParse(helper.RandomOneOfLength([]string{
					"d4feb120-9d7a-5d37-ab24-0c020e63210f",
					"35b453a5-3639-5afb-b99b-276fafccf913",
					"eec2b1cb-5864-4e20-907d-8b9635d9e840",
					"8991d0ac-2b43-5280-bb5c-aaa4ba3c6dd5",
					"8ae83e2e-1328-5f47-8a03-55dd690400a3",
					"d09e7c88-0a01-5c99-b7e4-692e9aea5ef0",
					"19f3d1a5-6066-5427-b30a-6b26c4337922",
				}))),
				Todo:     lo.ToPtr(helper.RandomOneOfLength([]string{"ติดตาม", "ส่งเอกสารเสนอราคา", "เจรจาต่อรอง"})),
				TodoAt:   lo.ToPtr(helper.RandomTime()),
				Type:     lo.ToPtr(helper.RandomDatatypesJSONFromSliceString([]string{"buyer", "seller"})),
				Interest: lo.ToPtr(helper.RandomDatatypesJSONFromSliceString([]string{"buy", "sell", "manage"})),
				Status:   lo.ToPtr(helper.RandomOneOfLength([]string{"booking", "offering", "survey", "negotiating", "contract", "new", "won_deal", "lost_deal", "contract_termination", "pre_booking"})),
				Tag:      lo.ToPtr(helper.RandomDatatypesJSONFromSliceString([]string{"condo", "sukhumvit", "sathorn", "silom", "ratchada", "ladprao", "ramkhamhaeng", "bangna", "bangkok", "swimming_pool", "gym", "garden", "playground", "security", "pet_friendly", "near_bts", "near_mrt", "near_airport", "near_expressway", "near_hospital", "near_school", "near_university", "near_shopping_mall", "near_department_store", "near_supermarket", "near_market", "near_park", "near_temple", "near_public_park", "near_public_transport", "near_public_pier", "near_public_boat", "near_public_bus", "near_public_train", "near_public_taxi", "near_public_tuk_tuk", "near_public_motorcycle", "near_public_bicycle", "near_public_car", "near_public_van", "near_public_truck"})),
				Gender:   lo.ToPtr(helper.RandomOneOfLength([]string{"MALE", "FEMALE"})),
			}
		}); err != nil {
			panic(err)
		}
		ch <- 1
	}()

	for i := 0; i < maxGoCount; i++ {
		<-ch
	}

	defer app.Close(context.Background())

}
