package v1

import (
	"go_base/controller"
	"go_base/domain"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

func RegisterRoutes(g echoswagger.ApiGroup, cfg *domain.Config) {
	g.GET("/ping", func(c echo.Context) error { return c.JSON(http.StatusOK, controller.Message{Message: "pong"}) }).
		AddResponse(http.StatusOK, "pong", controller.Message{Message: "pong"}, nil)
}
