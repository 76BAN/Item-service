package handler

import (
	"items_service/model"
	"items_service/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type ItemHandler struct {
	ItemService *service.ItemService
}

func NewitemHandler(itemService *service.ItemService) *ItemHandler {
	return &ItemHandler{
		ItemService: itemService,
	}
}

func (h *ItemHandler) CreateItem(c *echo.Context) error {
	var item model.Item
	userID, ok := c.Get("userID").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "UserID tidak valid",
		})
	}

	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"massage": "Request tidak valid" + err.Error(),
		})
	}

	createdItem, err := h.ItemService.CreateItem(&item, userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, createdItem)
}

func (h *ItemHandler) FindByID(c *echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "ID tidak valid",
		})
	}
	userID, ok := c.Get("userID").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "UserID tidak valid",
		})
	}
	category, err := h.ItemService.FindItemByID(id, userID)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Item tidak ditemukan",
		})
	}

	return c.JSON(http.StatusOK, category)
}

func (h *ItemHandler) FindItemByCategoryID(c *echo.Context) error {
	userID, ok := c.Get("userID").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "UserID tidak valid",
		})
	}
	idParam := c.Param("categoryId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "ID tidak valid",
		})
	}
	category, err := h.ItemService.FindItemByCategoryId(id, userID)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Item tidak ditemukan" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, category)
}

func (h *ItemHandler) DeleteItem(c *echo.Context) error {
	userID, ok := c.Get("userID").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "UserID tidak valid",
		})
	}
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "ID tidak valid",
		})
	}

	err = h.ItemService.DeleteItem(id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Item berhasil dihapus",
	})
}

func (h *ItemHandler) UpdateItem(c *echo.Context) error {
	userID, ok := c.Get("userID").(int64)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "UserID tidak valid",
		})
	}
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "ID tidak valid",
		})
	}

	var Item model.Item

	if err := c.Bind(&Item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Request tidak valid",
		})
	}

	Item.ID = id

	updatedItem, err := h.ItemService.UpdateItem(&Item, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, updatedItem)
}
