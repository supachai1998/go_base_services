package v1

import (
	"go_base/controller"
	"go_base/controller/middleware"
	"go_base/domain"
	"go_base/domain/permission"
	"net/http"

	"github.com/pangpanglabs/echoswagger/v2"
)

func RegisterRoutesProject(g echoswagger.ApiGroup, cfg *domain.Config) {
	handler := controller.ProjectHandler{Services: cfg.Services}
	auth := middleware.Auth(cfg.AdminAuthSecret, cfg.UserAuthSecret, cfg.CacheFunc)
	attach := middleware.Attach(cfg.Services.User.Get, cfg.Services.Staff.Get)
	verify := middleware.Verify(cfg.Services.User.Get, cfg.Services.Staff.Get)

	g.SetSecurity(domain.AuthHeaderKeyStaff).SetDescription("Project")
	restrict := middleware.RestrictPermissions(cfg.Services.Role.HasPermission)

	// GET /projects
	g.GET("", handler.Find, auth, attach, verify).
		AddParamQueryNested(domain.PaginationSwagger{}).
		AddResponse(http.StatusOK, "OK", domain.Pagination[domain.Project]{}, nil)

	// GET /projects/:id
	g.GET("/:id", handler.Get, auth, attach, verify).
		AddParamPath("", "id", "ID").
		AddResponse(http.StatusOK, "OK", domain.Project{}, nil)

	// POST /projects
	g.POST("", handler.Create, auth, attach, verify, restrict(permission.PROJECT_CREATE_ALL)).
		AddParamFormNested(domain.ProjectCreate{}).
		AddResponse(http.StatusCreated, "OK", nil, nil)

	// Update /projects/:id
	g.PUT("/:id", handler.Update, auth, attach, verify, restrict(permission.PROJECT_UPDATE_ALL)).
		AddParamPath("", "id", "ID").
		AddParamFormNested(domain.ProjectUpdate{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// DELETE /projects/:id
	g.DELETE("/:id", handler.Delete, auth, attach, verify, restrict(permission.PROJECT_DELETE_ALL)).
		AddParamPath("", "id", "ID").
		AddResponse(http.StatusNoContent, "OK", nil, nil)

}
