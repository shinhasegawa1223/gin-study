package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestPostItemReturnsCreated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	itemRepository := repositories.NewItemMemoryRepository(nil)
	itemService := services.NewItemService(itemRepository)
	itemController := controller.NewItemController(itemService)

	router := gin.New()
	router.POST("/items", itemController.Create)

	body := `{"name":"new product","price":500,"description":"sample description"}`
	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d, body=%s", w.Code, w.Body.String())
	}

	var resp struct {
		Data models.Item `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.ID == 0 {
		t.Fatalf("expected generated ID, got %#v", resp.Data)
	}

	if resp.Data.Name != "new product" || resp.Data.Price != 500 {
		t.Fatalf("unexpected item fields: %#v", resp.Data)
	}
}
