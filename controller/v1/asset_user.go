package v1

import (
	"go_base/controller"
	"go_base/controller/middleware"
	"go_base/domain"
	"net/http"

	"github.com/pangpanglabs/echoswagger/v2"
)

func RegisterRoutesAssetUser(g echoswagger.ApiGroup, cfg *domain.Config) {
	handler := controller.AssetHandler{Services: cfg.Services}
	auth := middleware.Auth(cfg.AdminAuthSecret, cfg.UserAuthSecret, cfg.CacheFunc)
	attach := middleware.Attach(cfg.Services.User.Get, cfg.Services.Staff.Get)
	verify := middleware.Verify(cfg.Services.User.Get, cfg.Services.Staff.Get)

	g.SetSecurity(domain.AuthHeaderKeyUser).SetDescription("Asset")

	// GET /assets
	g.GET("", handler.FindUser, auth, attach, verify).
		AddParamQueryNested(domain.PaginationSwagger{}).
		AddResponse(http.StatusOK, "OK", domain.Pagination[domain.Asset]{}, nil)

	// GET /assets/:id
	g.GET("/:id", handler.GetUser, auth, attach, verify).
		AddParamPath("", "id", "ID").
		AddResponse(http.StatusOK, "OK", domain.Asset{}, nil)

	// POST /assets
	g.POST("", handler.CreateUser, auth, attach, verify).
		AddParamFormNested(domain.AssetCreate{}).
		AddResponse(http.StatusCreated, "OK", nil, nil)

	// Update /assets/:id
	g.PUT("/:id", handler.UpdateUser, auth, attach, verify).
		AddParamPath("", "id", "ID").
		AddParamFormNested(domain.AssetUpdate{}).
		AddResponse(http.StatusOK, "OK", nil, nil)

	// DELETE /assets/:id
	g.DELETE("/:id", handler.DeleteUser, auth, attach, verify).
		AddParamPath("", "id", "ID").
		AddResponse(http.StatusNoContent, "OK", nil, nil)

}
