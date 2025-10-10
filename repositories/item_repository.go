package repositories

import (
	"errors"

	"gin-fleamarket/models"
)

type IItemRepository interface {
	FindAll() (*[]models.Item, error)
	FindbyId(itemId uint) (*models.Item, error)
	Create(item models.Item) (*models.Item, error)
}

type ItemMemoryRepository struct {
	items  []models.Item
	nextID uint
}

func NewItemMemoryRepository(items []models.Item) IItemRepository {
	repo := &ItemMemoryRepository{
		items:  make([]models.Item, len(items)),
		nextID: 1,
	}

	copy(repo.items, items)
	for _, item := range repo.items {
		if item.ID >= repo.nextID {
			repo.nextID = item.ID + 1
		}
	}

	return repo
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

func (r *ItemMemoryRepository) Create(item models.Item) (*models.Item, error) {
	if item.ID == 0 {
		item.ID = r.nextID
		r.nextID++
	} else {
		for _, existing := range r.items {
			if existing.ID == item.ID {
				return nil, errors.New("item already exists")
			}
		}
		if item.ID >= r.nextID {
			r.nextID = item.ID + 1
		}
	}

	r.items = append(r.items, item)
	return &r.items[len(r.items)-1], nil
}
