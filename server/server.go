package server

import (
	"context"
	"fmt"
	"go_base/configs"
	"go_base/controller"
	"go_base/controller/middleware"
	v1 "go_base/controller/v1"
	"go_base/domain"
	"go_base/logger"
	"go_base/validate"
	"go_base/xerror"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/pangpanglabs/echoswagger/v2"
	"go.uber.org/zap"
)

// CutstomValidator :
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate : Validate Data
func (cv *CustomValidator) Validate(i interface{}) error {
	return validate.Struct(i)
}

func Run(ctx context.Context) error {
	// Run starts the server.
	serverCtx, shutdown := context.WithCancel(ctx)

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		shutdown()
	}()

	app, err := CreateApp(ctx)
	if err != nil {
		return fmt.Errorf("CreateApp: %v", err)
	}
	defer app.Close(ctx)

	return run(serverCtx, app)
}

func panicRecoverConfig() echomw.RecoverConfig {
	cfg := echomw.DefaultRecoverConfig

	cfg.LogErrorFunc = func(c echo.Context, err error, stack []byte) error {
		return xerror.EInternal().SetDebugInfo("stack_trace", fmt.Sprintf("[PANIC RECOVER] %v %s", err, stack))
	}

	return cfg
}
func run(ctx context.Context, app *App) error {
	// Setup router
	echo.NotFoundHandler = func(c echo.Context) error {
		return xerror.E(xerror.ErrEndpointNotFound).SetStatusCode(xerror.ErrCodeNotFound)
	}
	// Load app configs
	cfg, err := configs.ParseConfig()
	if err != nil {
		return xerror.E(err).SetStatusCode(xerror.ErrCodeInternal)
	}
	apiV1 := "/api/v1"

	ewg := echoswagger.New(echo.New(), apiV1+"/doc", &echoswagger.Info{
		Title:          cfg.SwaggerTitle,
		Description:    cfg.SwaggerDescription,
		Version:        cfg.Version,
		TermsOfService: fmt.Sprintf("%s/terms/", cfg.BaseUrl),
		Contact: &echoswagger.Contact{
			Email: cfg.SwaggerContact.Email,
			Name:  cfg.SwaggerContact.Name,
			URL:   cfg.SwaggerContact.URL,
		},
		License: &echoswagger.License{
			Name: cfg.SwaggerLicense.Name,
			URL:  cfg.SwaggerLicense.URL,
		}}).
		AddSecurityAPIKey(domain.AuthHeaderKeyStaff, "staff", echoswagger.SecurityInHeader).
		AddSecurityAPIKey(domain.AuthHeaderKeyUser, "user", echoswagger.SecurityInHeader)
	e := ewg.Echo()
	var se echoswagger.ApiRoot
	if app.Cfg.ENV != "production" {
		se = echoswagger.NewNop(e)
	} else {
		se = ewg
	}
	e = se.Echo()

	gV1 := ewg.Group("v1", apiV1)
	e.HideBanner = true
	e.HTTPErrorHandler = controller.ErrorHandler(app.Cfg.DebugMode)

	loggerWriter := zap.NewStdLog(logger.L().Desugar()).Writer()
	e.Logger.SetOutput(loggerWriter)

	e.Pre(echomw.RemoveTrailingSlash())
	e.Use(
		echomw.RecoverWithConfig(panicRecoverConfig()),
		middleware.RequestLogger(app.Cfg.ENV, app.Cfg.ECHODATA.REQ, app.Cfg.ECHODATA.RES),
		middleware.BodyDumpWithConfig(app.Cfg.ENV),
		middleware.CorrelationID(),
		echomw.SecureWithConfig(echomw.SecureConfig{
			ContentTypeNosniff: "nosniff",
			HSTSMaxAge:         63072000, // 2 years
		}),
		echomw.TimeoutWithConfig(echomw.TimeoutConfig{
			Timeout: app.Cfg.Server.RequestTimeout,
		}),
	)
	e.Static("/", "public")
	// Middleware bind to echo instance
	// logger.WriteLog()
	// e.Use(echomw.LoggerWithConfig(echomw.LoggerConfig{
	// 	Format: "${time_rfc3339_nano} ${remote_ip}::${method}${path}, status=${status}, latency=${latency_human}, error=${error}\n",
	// 	Output: io.MultiWriter(os.Stdout, logger.LogFile),
	// }))
	e.Validator = &CustomValidator{Validator: validate.New()}
	healthCheckRes := controller.Success{Success: true}
	gV1.GET("/health", func(c echo.Context) error { return c.JSON(http.StatusOK, healthCheckRes) }).AddResponse(http.StatusOK, "health check", healthCheckRes, nil)

	v1.RegisterRoutes(gV1, &domain.Config{
		Services:  app.Services,
		CacheFunc: app.Redis.GetStringValue,
	})

	// staff
	groupStaff := ewg.Group("staff", apiV1+"/staffs")
	v1.RegisterRoutesStaff(groupStaff, &domain.Config{
		Services:        app.Services,
		CacheFunc:       app.Redis.GetStringValue,
		AdminAuthSecret: cfg.AdminAuth.JWTSecret,
		UserAuthSecret:  cfg.UserAuth.JWTSecret,
	})
	// staff me
	groupStaffMe := ewg.Group("staff_me", apiV1+"/me")
	v1.RegisterRoutesStaffMe(groupStaffMe, &domain.Config{
		Services:        app.Services,
		CacheFunc:       app.Redis.GetStringValue,
		AdminAuthSecret: cfg.AdminAuth.JWTSecret,
		UserAuthSecret:  cfg.UserAuth.JWTSecret,
	})

	// role
	groupRole := ewg.Group("role", apiV1+"/roles")
	v1.RegisterRoutesRole(groupRole, &domain.Config{
		Services:        app.Services,
		CacheFunc:       app.Redis.GetStringValue,
		AdminAuthSecret: cfg.AdminAuth.JWTSecret,
		UserAuthSecret:  cfg.UserAuth.JWTSecret,
	})

	// user
	groupUser := ewg.Group("user", apiV1+"/users")
	v1.RegisterRoutesUser(groupUser, &domain.Config{
		Services:        app.Services,
		CacheFunc:       app.Redis.GetStringValue,
		AdminAuthSecret: cfg.AdminAuth.JWTSecret,
		UserAuthSecret:  cfg.UserAuth.JWTSecret,
	})

	// developer
	groupDeveloper := ewg.Group("developer", apiV1+"/developers")
	v1.RegisterRoutesDeveloper(groupDeveloper, &domain.Config{
		Services:        app.Services,
		CacheFunc:       app.Redis.GetStringValue,
		AdminAuthSecret: cfg.AdminAuth.JWTSecret,
		UserAuthSecret:  cfg.UserAuth.JWTSecret,
	})

	// project
	groupProject := ewg.Group("project", apiV1+"/projects")
	v1.RegisterRoutesProject(groupProject, &domain.Config{
		Services:        app.Services,
		CacheFunc:       app.Redis.GetStringValue,
		AdminAuthSecret: cfg.AdminAuth.JWTSecret,
		UserAuthSecret:  cfg.UserAuth.JWTSecret,
	})

	// asset
	groupAsset := ewg.Group("asset", apiV1+"/assets")
	v1.RegisterRoutesAsset(groupAsset, &domain.Config{
		Services:        app.Services,
		CacheFunc:       app.Redis.GetStringValue,
		AdminAuthSecret: cfg.AdminAuth.JWTSecret,
		UserAuthSecret:  cfg.UserAuth.JWTSecret,
	})

	// asset user
	groupAssetUser := ewg.Group("asset user", apiV1+"/assets/user")
	v1.RegisterRoutesAssetUser(groupAssetUser, &domain.Config{
		Services:        app.Services,
		CacheFunc:       app.Redis.GetStringValue,
		AdminAuthSecret: cfg.AdminAuth.JWTSecret,
		UserAuthSecret:  cfg.UserAuth.JWTSecret,
	})

	errCh := make(chan error)

	// Run the server
	go func() {
		logger.L().Infof("⇨ env:%s || version:%s", app.Cfg.ENV, app.Cfg.Version)
		if err := e.Start(fmt.Sprintf(":%d", app.Cfg.Server.Port)); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Wait for server context to be stopped
	select {
	case err := <-errCh:
		logger.L().Errorf("Server error: %s\n", err)
		break
	case <-ctx.Done():
		break
	}

	logger.L().Info("Shutting down a server")

	// Shutdown signal with grace period of 30 seconds
	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), app.Cfg.Server.ShutdownTimeout)
	shutdownErrCh := make(chan error)

	// Trigger graceful shutdown
	go func() {
		if err := e.Shutdown(shutdownCtx); err != nil {
			shutdownErrCh <- err
		}

		cancelShutdownCtx()
	}()

	select {
	case err := <-shutdownErrCh:
		logger.L().Errorf("Graceful shutdown error: %s\n", err)
	case <-shutdownCtx.Done():
		if err := shutdownCtx.Err(); err == context.DeadlineExceeded {
			return fmt.Errorf("graceful shutdown timeout: %w", err)
		}
		logger.L().Info("Graceful shutdown success")
	}
	return nil

}
