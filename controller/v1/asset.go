package v1

import (
	"go_base/controller"
	"go_base/controller/middleware"
	"go_base/domain"
	"go_base/domain/permission"
	"net/http"

	"github.com/pangpanglabs/echoswagger/v2"
)

func RegisterRoutesAsset(g echoswagger.ApiGroup, cfg *domain.Config) {
	handler := controller.AssetHandler{Services: cfg.Services}
	auth := middleware.Auth(cfg.AdminAuthSecret, cfg.UserAuthSecret, cfg.CacheFunc)
	attach := middleware.Attach(cfg.Services.User.Get, cfg.Services.Staff.Get)
	verify := middleware.Verify(cfg.Services.User.Get, cfg.Services.Staff.Get)

	g.SetSecurity(domain.AuthHeaderKeyStaff).SetDescription("Asset")
	restrict := middleware.RestrictPermissions(cfg.Services.Role.HasPermission)

	// GET /assets
	g.GET("", handler.Find, auth, attach, verify, restrict(permission.ASSET_VIEW_ALL)).
		AddParamQueryNested(domain.PaginationSwagger{}).
		AddResponse(http.StatusOK, "OK", domain.Pagination[domain.Asset]{}, nil)

	// GET /assets/:id
	g.GET("/:id", handler.Get, auth, attach, verify, restrict(permission.ASSET_VIEW_ALL)).
		AddParamPath("", "id", "ID").
		AddResponse(http.StatusOK, "OK", domain.Asset{}, nil)

	// POST /assets
	g.POST("", handler.Create, auth, attach, verify, restrict(permission.ASSET_CREATE_ALL)).
		AddParamFormNested(domain.AssetCreate{}).
		AddResponse(http.StatusCreated, "OK", nil, nil)

	// Update /assets/:id
	g.PUT("/:id", handler.Update, auth, attach, verify, restrict(permission.ASSET_UPDATE_ALL)).
		AddParamPath("", "id", "ID").
		AddParamFormNested(domain.AssetUpdate{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// DELETE /assets/:id
	g.DELETE("/:id", handler.Delete, auth, attach, verify, restrict(permission.ASSET_DELETE_ALL)).
		AddParamPath("", "id", "ID").
		AddResponse(http.StatusNoContent, "OK", nil, nil)

}
