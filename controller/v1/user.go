package v1

import (
	"go_base/controller"
	"go_base/controller/middleware"
	"go_base/domain"
	"go_base/domain/permission"
	"net/http"

	"github.com/pangpanglabs/echoswagger/v2"
)

func RegisterRoutesUser(g echoswagger.ApiGroup, cfg *domain.Config) {
	handler := controller.UserHandler{Services: cfg.Services}
	auth := middleware.Auth(cfg.AdminAuthSecret, cfg.UserAuthSecret, cfg.CacheFunc)
	attach := middleware.Attach(cfg.Services.User.Get, cfg.Services.Staff.Get)
	verify := middleware.Verify(cfg.Services.User.Get, cfg.Services.Staff.Get)

	g.SetSecurity(domain.AuthHeaderKeyUser).SetDescription("User")
	restrict := middleware.RestrictPermissions(cfg.Services.Role.HasPermission)

	// GET /users
	g.GET("", handler.Find, auth, attach, verify, restrict(permission.USER_VIEW_ALL)).
		AddParamQueryNested(domain.PaginationSwagger{}).
		AddResponse(http.StatusOK, "OK", domain.Pagination[domain.User]{}, nil)

	// POST /users/token
	g.POST("/token", handler.GetToken).
		AddParamFormNested(domain.UserGetToken{}).
		AddResponse(http.StatusOK, "OK", domain.UserGetTokenResponse{}, nil)

	// POST /users
	g.POST("", handler.Create).
		AddParamFormNested(domain.UserCreate{}).
		AddResponse(http.StatusOK, "OK", domain.User{}, nil)

	// POST /users/login
	g.POST("/login", handler.LoginWithEmailPassword).
		AddParamFormNested(domain.UserLogin{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// POST /users/unlock
	g.POST("/unlock", handler.Unlock, auth, attach, verify, restrict(permission.USER_UNLOCK_ALL)).
		AddParamFormNested(domain.UserUnlock{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// DELETE /users/:id
	g.DELETE("/:id", handler.Delete, auth, attach, verify, restrict(permission.USER_DELETE_ALL)).
		AddParamPath("", "id", "user id").
		AddResponse(http.StatusOK, "OK", nil, nil)

	// UPDATE /users/:id
	g.PUT("/:id", handler.Update, auth, attach, verify, restrict(permission.USER_UPDATE_ALL)).
		SetSecurity(domain.AuthHeaderKeyUser).
		AddParamFormNested(domain.UserUpdate{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// POST /users/verify
	g.POST("/verify", handler.VerifyToken).
		AddParamFormNested(domain.UserVerifyToken{}).
		AddResponse(http.StatusOK, "OK", domain.UserVerifyTokenResponse{}, nil)

	// Update Password /users/password
	g.PUT("/me/password", handler.UpdatePassword, auth, attach).
		AddParamFormNested(domain.UserUpdatePassword{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// Update Me /users/me
	g.PUT("/me", handler.UpdateMe, auth, attach).
		AddParamFormNested(domain.UserUpdate{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// Get Me /users/me
	g.GET("/me", handler.GetMe, auth, attach).
		AddResponse(http.StatusOK, "OK", domain.UserMe{}, nil)

	// Get log me /users/me/log
	g.GET("/me/log", handler.GetLogMe, auth, attach, verify).
		AddParamQueryNested(domain.PaginationSwagger{}).
		AddResponse(http.StatusOK, "OK", domain.Pagination[domain.Logs[domain.User]]{}, nil)

	// Delete /users/ids
	g.DELETE("/ids", handler.DeleteIds, auth, attach, verify, restrict(permission.USER_DELETE_ALL)).
		AddParamFormNested(domain.Ids{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

}
