package v1

import (
	"go_base/controller"
	"go_base/controller/middleware"
	"go_base/domain"
	"go_base/domain/permission"
	"net/http"

	"github.com/pangpanglabs/echoswagger/v2"
)

func RegisterRoutesRole(g echoswagger.ApiGroup, cfg *domain.Config) {
	h := controller.RoleHandler{Services: cfg.Services}
	auth := middleware.Auth(cfg.AdminAuthSecret, cfg.UserAuthSecret, cfg.CacheFunc)
	attach := middleware.Attach(cfg.Services.User.Get, cfg.Services.Staff.Get)
	verify := middleware.Verify(cfg.Services.User.Get, cfg.Services.Staff.Get)

	g.SetSecurity(domain.AuthHeaderKeyStaff).SetDescription("Staff")
	restrict := middleware.RestrictPermissions(cfg.Services.Role.HasPermission)
	// Get role /roles
	g.GET("", h.Find, auth, attach, verify).
		AddParamQueryNested(domain.PaginationSwagger{}).
		AddResponse(http.StatusOK, "OK", domain.Pagination[domain.Role]{}, nil)

	// Create role /roles
	g.POST("", h.Create, auth, attach, verify, restrict(permission.ROLE_CREATE)).
		AddParamBody(domain.RoleSwaggerCreate{}, "body", "", true).
		AddResponse(http.StatusCreated, "OK", nil, nil)

	// Update role /roles/:id
	g.PUT("/:id", h.Update, auth, attach, verify, restrict(permission.ROLE_UPDATE)).
		AddParamPath("", "id", "ID").
		AddParamBody(domain.RoleUpdate{}, "body", "", true).
		AddResponse(http.StatusOK, "OK", nil, nil)

}
