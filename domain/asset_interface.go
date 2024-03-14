package domain

import "github.com/labstack/echo/v4"

type IAssetService[T Asset, U AssetUpdate, C AssetCreate] interface {
	CreateScope(ctx echo.Context, c *AssetCreate) error
}
