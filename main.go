package main

import (
	"log"

	"gin-fleamarket/controller"
	"gin-fleamarket/infra"
	"gin-fleamarket/models"
	"gin-fleamarket/repositories"
	"gin-fleamarket/services"

	"github.com/gin-gonic/gin"
)

func main() {
	infra.Initialize()

	db := infra.SetupDB()
	if err := db.AutoMigrate(&models.Item{}); err != nil {
		log.Fatalf("failed to migrate items table: %v", err)
	}

	itemRepository := repositories.NewItemRepository(db)
	itemService := services.NewItemService(itemRepository)
	itemController := controller.NewItemController(itemService)

	router := gin.Default()
	router.GET("/items", itemController.FindAll)
	router.GET("/items/:id", itemController.FindbyId)
	router.POST("/items", itemController.Create)
	router.PUT("/items/:id", itemController.Update)
	router.DELETE("/items/:id", itemController.Delete)

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
