package v1

import (
	"go_base/controller"
	"go_base/controller/middleware"
	"go_base/domain"
	"go_base/domain/permission"
	"net/http"

	"github.com/pangpanglabs/echoswagger/v2"
)

func RegisterRoutesStaffMe(g echoswagger.ApiGroup, cfg *domain.Config) {
	handler := controller.StaffMeHandler{Services: cfg.Services}
	auth := middleware.Auth(cfg.AdminAuthSecret, cfg.UserAuthSecret, cfg.CacheFunc)
	attach := middleware.Attach(cfg.Services.User.Get, cfg.Services.Staff.Get)
	verify := middleware.Verify(cfg.Services.User.Get, cfg.Services.Staff.Get)

	g.SetSecurity(domain.AuthHeaderKeyStaff).SetDescription("Staff")
	restrict := middleware.RestrictPermissions(cfg.Services.Role.HasPermission)
	// Get log me delete /staff/log
	g.GET("/log", handler.GetLogMeDelete, auth, attach, verify, restrict(permission.STAFF_ME_LOG_SELF)).
		AddParamQueryNested(domain.PaginationSwagger{}).
		AddResponse(http.StatusOK, "OK", domain.Pagination[domain.Logs[domain.Staff]]{}, nil)
	// Update Password /staff/password
	g.PUT("/password", handler.UpdatePassword, auth, attach).
		AddParamFormNested(domain.StaffUpdatePassword{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// Get me /staff/me
	g.GET("", handler.GetMe, auth, attach).
		AddResponse(http.StatusOK, "OK", domain.StaffMe{}, nil)
}
