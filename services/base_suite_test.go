package services_test

import (
	"context"
	"fmt"
	"go_base/database"
	"go_base/domain"
	"go_base/server"
	"go_base/storage"
	"go_base/validate"
	"net/http"
	"testing"
	"time"

	"go_base/configs"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// Integration tests
type SetupSuite interface {
	SetupSuite()
	TearDownSuite()
}

type App struct {
	Services *domain.AllServices
	Stores   *database.AllStores
	Storages *storage.AllStorage
	Cfg      *configs.Config
	DB       *gorm.DB
	Redis    *storage.Cache
}

func (app *App) Close(ctx context.Context) {
	// Handling close connection
	if app.DB == nil || app.Redis == nil {
		panic("cannot close db or redis connection")
	}

	db, err := app.DB.DB()
	if err != nil {
		panic("cannot get db connection")
	}
	_ = db.Close()
	// redis flush all
	_ = app.Redis.Client.Close()
}

type UnitTestSuite struct {
	suite.Suite
	server  *App
	service *domain.AllServices
	ctx     echo.Context
}

func (its *UnitTestSuite) SetupSuite() {
	s, err := server.CreateAppForTest(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	its.server = &App{
		Services: s.Services,
		Stores:   s.Stores,
		Storages: s.Storages,
		Cfg:      s.Cfg,
		DB:       s.DB,
		Redis:    s.Redis,
	}
	its.ctx = *MockEchoContext()

	its.service = s.Services

	_ = s.Redis.Client.FlushAll(its.ctx.Request().Context()).Err()

}

func (s *UnitTestSuite) TearDownSuite() {
	time.Sleep(1 * time.Second)
	s.server.Close(context.Background())
}
func TestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func MockEchoContext() *echo.Context {
	e := echo.New()
	e.Validator = &server.CustomValidator{Validator: validate.New()}
	return lo.ToPtr(e.NewContext(
		&http.Request{Header: http.Header{"Content-Type": []string{"application/json"}}},
		echo.NewResponse(nil, nil),
	))
}
