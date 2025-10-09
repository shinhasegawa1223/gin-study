package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-fleamarket/controller"
	"gin-fleamarket/models"
	"gin-fleamarket/repositories"
	"gin-fleamarket/services"

	"github.com/gin-gonic/gin"
)

func TestGetItemByIDReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	items := []models.Item{
		{ID: 1, Name: "product1", Price: 100, Description: "description1", SoldOut: false},
		{ID: 2, Name: "product2", Price: 200, Description: "description2", SoldOut: false},
	}

	itemRepository := repositories.NewItemMemoryRepository(items)
	itemService := services.NewItemService(itemRepository)
	itemController := controller.NewItemController(itemService)

	router := gin.New()
	router.GET("/items/:id", itemController.FindbyId)

	req := httptest.NewRequest(http.MethodGet, "/items/2", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", w.Code, w.Body.String())
	}
}
