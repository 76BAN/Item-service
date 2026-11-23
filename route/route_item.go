package routes

import (
	"items_service/handler"
	middleware "items_service/midleware"

	"github.com/labstack/echo/v5"
)

func Init(
	e *echo.Echo,
	itemHandler *handler.ItemHandler,
	midleware *middleware.JWTMiddleware,

) {
	e.GET("/groups/category/item/:id", itemHandler.FindByID, midleware.Authenticate)
	e.GET("/groups/category/itembyCategory/:categoryId", itemHandler.FindItemByCategoryID, midleware.Authenticate)
	e.POST("/groups/category/item", itemHandler.CreateItem, midleware.Authenticate)
	e.DELETE("/groups/category/item/:id", itemHandler.DeleteItem, midleware.Authenticate)
	e.PUT("/groups/category/item/:id", itemHandler.UpdateItem, midleware.Authenticate)
 
}
