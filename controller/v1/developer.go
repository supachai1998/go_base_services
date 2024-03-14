package v1

import (
	"go_base/controller"
	"go_base/controller/middleware"
	"go_base/domain"
	"go_base/domain/permission"
	"net/http"

	"github.com/pangpanglabs/echoswagger/v2"
)

func RegisterRoutesDeveloper(g echoswagger.ApiGroup, cfg *domain.Config) {
	handler := controller.DeveloperHandler{Services: cfg.Services}
	auth := middleware.Auth(cfg.AdminAuthSecret, cfg.UserAuthSecret, cfg.CacheFunc)
	attach := middleware.Attach(cfg.Services.User.Get, cfg.Services.Staff.Get)
	verify := middleware.Verify(cfg.Services.User.Get, cfg.Services.Staff.Get)

	g.SetSecurity(domain.AuthHeaderKeyStaff).SetDescription("Developer")
	restrict := middleware.RestrictPermissions(cfg.Services.Role.HasPermission)

	// GET /developers
	g.GET("", handler.Find, auth, attach, verify).
		AddParamQueryNested(domain.PaginationSwagger{}).
		AddResponse(http.StatusOK, "OK", domain.Pagination[domain.Developer]{}, nil)

	// GET /developers/:id
	g.GET("/:id", handler.Get, auth, attach, verify).
		AddParamPath("", "id", "ID").
		AddResponse(http.StatusOK, "OK", domain.Developer{}, nil)

	// POST /developers
	g.POST("", handler.Create, auth, attach, verify, restrict(permission.DEVELOPER_CREATE_ALL)).
		AddParamFormNested(domain.DeveloperCreate{}).
		AddResponse(http.StatusCreated, "OK", nil, nil)

	// Update /developers/:id
	g.PUT("/:id", handler.Update, auth, attach, verify, restrict(permission.DEVELOPER_UPDATE_ALL)).
		AddParamPath("", "id", "ID").
		AddParamFormNested(domain.DeveloperUpdate{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// DELETE /developers/:id
	g.DELETE("/:id", handler.Delete, auth, attach, verify, restrict(permission.DEVELOPER_DELETE_ALL)).
		AddParamPath("", "id", "ID").
		AddResponse(http.StatusNoContent, "OK", nil, nil)

}
