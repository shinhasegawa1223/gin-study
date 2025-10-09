package repositories

import (
	"errors"
	"gin-fleamarket/models"
)

type IItemRepository interface {
	FindAll() (*[]models.Item, error)
	FindbyId(itemId uint) (*models.Item, error)
}

type ItemMemoryRepository struct {
	items []models.Item
}

func NewItemMemoryRepository(items []models.Item) IItemRepository {
	return &ItemMemoryRepository{items: items}
}

func (r *ItemMemoryRepository) FindAll() (*[]models.Item, error) {
	return &r.items, nil
}

func (r *ItemMemoryRepository) FindbyId(itemId uint) (*models.Item, error) {
	for i := range r.items {
		if r.items[i].ID == itemId {
			return &r.items[i], nil
		}
	}
	return nil, errors.New("item not found")

}
