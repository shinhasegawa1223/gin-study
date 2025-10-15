package repositories

import (
	"errors"
	"gin-fleamarket/models"
)

type IItemRepository interface {
	FindAll() (*[]models.Item, error)
	FindbyId(itemId uint) (*models.Item, error)
	Create(NewItem models.Item) (*models.Item, error)
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


func(r *ItemMemoryRepository) Create(NewItem models.Item) (*models.Item, error){
	NewItem.ID = uint(len(r.items) + 1)
	r.items = append(r.items, NewItem)
	return &NewItem, nil
}