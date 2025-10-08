package main

import (
	"gin-fleamarket/controller"
	"gin-fleamarket/models"
	"gin-fleamarket/repositories"
	"gin-fleamarket/services"

	"github.com/gin-gonic/gin"
)

func main() {
	items := []models.Item{
		{ID: 1, Name: "product1", Price: 100, Description: "description1", SoldOut: false},
		{ID: 2, Name: "2product2", Price: 102220, Description: "2", SoldOut: true},
		{ID: 3, Name: "product3", Price: 333333, Description: "3", SoldOut: false},
	}
	itemRepository := repositories.NewItemMemoryRepository(items)
	itemService := services.NewItemService(itemRepository)
	itemController := controller.NewItemController(itemService)

	router := gin.Default()
	router.GET("/items", itemController.FindAll)
	router.Run("localhost:8080")
}
