package server

import (
	"context"
	"flag"
	"fmt"
	"go_base/configs"
	"go_base/database"
	"go_base/domain"
	myLogger "go_base/logger"
	"go_base/services"
	"go_base/services/auth"
	"go_base/storage"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

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
		myLogger.L().Error("cannot close db or redis connection")
	}

	db, err := app.DB.DB()
	if err != nil {
		myLogger.L().Error("cannot get db connection")
	}
	_ = db.Close()

	_ = app.Redis.Client.Close()
}

func CreateApp(ctx context.Context) (*App, error) {
	loadDotenv := flag.Bool("dotenv", false, "Load app config from .env file.")

	flag.Parse()

	myLogger.MustInit()

	if *loadDotenv {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("load .env failed: %v", err)
		}
	}
	// Load app configs
	cfg, err := configs.ParseConfig()
	if err != nil {
		return nil, fmt.Errorf("parse config failed: %v", err)
	}
	app, err := createApp(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("createApp: %v", err)
	}

	return app, nil
}

func CreateAppForTest(ctx context.Context) (*App, error) {
	loadDotenv := flag.Bool("dotenv", false, "Load app config from .env file.")

	flag.Parse()

	if *loadDotenv {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("load .env failed: %v", err)
		}
	}
	// Load app configs
	cfg, err := configs.ParseConfig(configs.ModeTest)
	if err != nil {
		return nil, fmt.Errorf("parse config failed: %v", err)
	}
	app, err := createApp(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("createApp: %v", err)
	}

	return app, nil
}

func createApp(ctx context.Context, cfg *configs.Config) (*App, error) {

	logFile := logFile("app")
	defer logFile.Close()
	if cfg.PrettyLog {
		log, err := myLogger.NewPretty(logFile)
		if err != nil {
			myLogger.L().Warn("cannot init pretty myLogger: %v; fallback to default", err)
		} else {
			myLogger.ReplaceGlobals(log)
		}

	} else {
		log.SetOutput(logFile)
	}

	// Init postgresql
	dsn := cfg.DB.GetURI()
	logInfo := logger.Info
	if !cfg.DB.DBDEBUG {
		logInfo = logger.Silent
	}
	postgresql, err := storage.NewPostgresClient(dsn, &gorm.Config{
		Logger:         logger.Default.LogMode(logInfo),
		TranslateError: false,
	}, cfg)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql: %v", err)
	}
	// Init redis
	redis, err := storage.NewRedisClient(ctx, cfg.Redis.GetOptions(), cfg.Redis.GetTLSServer())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	adminAuthCfg := auth.AuthConfig(cfg.AdminAuth)
	userAuthCfg := auth.AuthConfig(cfg.UserAuth)

	allStorage := &storage.AllStorage{
		DB:    postgresql.Client,
		Cache: redis,
	}

	// store
	store := database.NewStore(postgresql.Client, redis, allStorage)
	stores := &database.AllStores{
		Base:      store,
		Staff:     database.NewStaffStore(postgresql.Client, allStorage),
		Auth:      database.NewAuthStore(postgresql.Client, allStorage),
		Role:      database.NewRoleStore(postgresql.Client, allStorage),
		User:      database.NewUserStore(postgresql.Client, allStorage),
		Developer: database.NewBaseStore[domain.Developer, domain.DeveloperUpdate, domain.DeveloperCreate](postgresql.Client, &database.BaseStoreConfig{WriteChangelog: true, CacheExpire: 10 * time.Minute}, allStorage),
		Project:   database.NewBaseStore[domain.Project, domain.ProjectUpdate, domain.ProjectCreate](postgresql.Client, &database.BaseStoreConfig{WriteChangelog: true, CacheExpire: 10 * time.Minute}, allStorage),
		Asset:     database.NewBaseStore[domain.Asset, domain.AssetUpdate, domain.AssetCreate](postgresql.Client, &database.BaseStoreConfig{WriteChangelog: true, CacheExpire: 10 * time.Minute}, allStorage),
	}

	// all services
	allServices := &domain.AllServices{}
	allServices.Role = services.NewRoleService(store, stores.Role, allServices, redis)
	allServices.Staff = services.NewStaffService(store, stores.Staff, allServices, redis, &adminAuthCfg)
	allServices.AuthAdmin = services.NewAuthAdminService(stores.Auth, allServices)
	allServices.AuthUser = services.NewAuthUserService(stores.Auth, allServices)
	allServices.User = services.NewUserService(store, stores.User, allServices, redis, &userAuthCfg)
	allServices.IDeveloper = services.NewBaseService(store, stores.Developer, allServices, redis)
	allServices.IProject = services.NewBaseService(store, stores.Project, allServices, redis)
	allServices.IAsset = services.NewBaseService(store, stores.Asset, allServices, redis)
	allServices.Asset = services.NewAssetService(store, stores.Asset, allServices, redis)
	return &App{
		Cfg:      cfg,
		DB:       postgresql.Client,
		Stores:   stores,
		Redis:    redis,
		Services: allServices,
		Storages: &storage.AllStorage{
			DB:    postgresql.Client,
			Cache: redis,
		},
	}, nil
}

func logGorm(debug, slient bool) logger.Interface {
	loggerZap := zapgorm2.New(zap.L())
	LogLevel := logger.Info
	if !debug {
		LogLevel = logger.Error
	}
	if slient {
		LogLevel = logger.Silent
	}
	loggerZap.LogLevel = LogLevel
	loggerZap.SetAsDefault()
	return loggerZap
}

func logFile(name string) *os.File {
	// write log
	t := time.Now().Local()
	fileName := t.Format("2006-01-02") + ".log"
	folder := configs.Root + "logger/log" + "/" + name
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, 0755)
	}
	pathJoin := filepath.Join(folder, fileName)
	if err := os.MkdirAll(filepath.Dir(pathJoin), 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}
	file, err := os.OpenFile(pathJoin, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	return file
}
