package main

import (
	"context"
	"go_base/domain"
	"go_base/server"
	"go_base/storage"
)

func main() {
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

	defer app.Close(context.Background())

}
