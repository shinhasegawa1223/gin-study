package controller

import (
	"net/http"
	"strconv"

	"gin-fleamarket/dto"
	"gin-fleamarket/services"

	"github.com/gin-gonic/gin"
)

type IItemController interface {
	FindAll(ctx *gin.Context)
	FindbyId(ctx *gin.Context)
	Create(ctx *gin.Context)
}

type itemController struct {
	service services.IItemService
}

func NewItemController(service services.IItemService) IItemController {
	return &itemController{service: service}
}

func (c *itemController) FindAll(ctx *gin.Context) {
	items, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": items})
}

func (c *itemController) FindbyId(ctx *gin.Context) {
	itemID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := c.service.FindbyId(uint(itemID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": item})
}

func (c *itemController) Create(ctx *gin.Context) {
	var input dto.CreateItemInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newItem, err := c.service.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": newItem})
}
