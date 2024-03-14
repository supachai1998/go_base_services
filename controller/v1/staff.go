package v1

import (
	"go_base/controller"
	"go_base/controller/middleware"
	"go_base/domain"
	"go_base/domain/permission"
	"net/http"

	"github.com/pangpanglabs/echoswagger/v2"
)

func RegisterRoutesStaff(g echoswagger.ApiGroup, cfg *domain.Config) {
	handler := controller.StaffHandler{Services: cfg.Services}
	auth := middleware.Auth(cfg.AdminAuthSecret, cfg.UserAuthSecret, cfg.CacheFunc)
	attach := middleware.Attach(cfg.Services.User.Get, cfg.Services.Staff.Get)
	verify := middleware.Verify(cfg.Services.User.Get, cfg.Services.Staff.Get)

	g.SetSecurity(domain.AuthHeaderKeyStaff).SetDescription("Staff")
	restrict := middleware.RestrictPermissions(cfg.Services.Role.HasPermission)

	// GET /staff
	g.GET("", handler.Find, auth, attach, verify, restrict(permission.STAFF_VIEW_ALL)).
		AddParamQueryNested(domain.PaginationSwagger{}).
		AddResponse(http.StatusOK, "OK", domain.Pagination[domain.Staff]{}, nil)

	// POST /staff/token
	g.POST("/token", handler.GetToken).
		AddParamFormNested(domain.StaffGetToken{}).
		AddResponse(http.StatusOK, "OK", domain.StaffGetTokenResponse{}, nil)

	// POST /staff
	g.POST("", handler.Create).
		AddParamFormNested(domain.StaffCreate{}).
		AddResponse(http.StatusOK, "OK", domain.Staff{}, nil)

	// POST /staff/login
	g.POST("/login", handler.LoginWithEmailPassword).
		AddParamFormNested(domain.StaffLogin{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// POST /staff/unlock
	g.POST("/unlock", handler.Unlock, auth, attach, verify, restrict(permission.STAFF_UNLOCK_ALL)).
		AddParamFormNested(domain.StaffUnlock{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// DELETE /staff/:id
	g.DELETE("/:id", handler.Delete, auth, attach, verify, restrict(permission.STAFF_DELETE_ALL)).
		AddParamPath("", "id", "staff id").
		AddResponse(http.StatusOK, "OK", nil, nil)

	// UPDATE /staff/:id
	g.PUT("/:id", handler.Update, auth, attach, verify, restrict(permission.STAFF_UPDATE_ALL)).
		SetSecurity(domain.AuthHeaderKeyStaff).
		AddParamFormNested(domain.StaffUpdate{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// POST /staff/verify
	g.POST("/verify", handler.VerifyToken).
		AddParamFormNested(domain.StaffVerifyToken{}).
		AddResponse(http.StatusOK, "OK", domain.StaffVerifyTokenResponse{}, nil)

	// GET /staff/log
	g.GET("/log/:id", handler.GetLog, auth, attach, verify, restrict(permission.STAFF_VIEW_ALL)).
		AddParamQueryNested(domain.PaginationSwagger{}).
		AddResponse(http.StatusOK, "OK", domain.Pagination[domain.Logs[domain.Staff]]{}, nil)

}
