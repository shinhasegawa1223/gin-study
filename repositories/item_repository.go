package repositories

import (
	"errors"
	"fmt"
	"sync"

	"gin-fleamarket/models"

	"gorm.io/gorm"
)

var (
	// ErrItemNotFound is returned when the requested item does not exist.
	ErrItemNotFound = errors.New("item not found")
)

type IItemRepository interface {
	FindAll() (*[]models.Item, error)
	FindbyId(itemID uint) (*models.Item, error)
	Create(item models.Item) (*models.Item, error)
	Update(item models.Item) (*models.Item, error)
	Delete(itemID uint) error
}

// itemRepository persists items in the backing database.
type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) IItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) FindAll() (*[]models.Item, error) {
	var items []models.Item
	if err := r.db.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("fetch items: %w", err)
	}
	return &items, nil
}

func (r *itemRepository) FindbyId(itemID uint) (*models.Item, error) {
	var item models.Item
	if err := r.db.First(&item, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("fetch item id=%d: %w", itemID, err)
	}
	return &item, nil
}

func (r *itemRepository) Create(item models.Item) (*models.Item, error) {
	if err := r.db.Create(&item).Error; err != nil {
		return nil, fmt.Errorf("create item: %w", err)
	}
	return &item, nil
}

func (r *itemRepository) Update(item models.Item) (*models.Item, error) {
	var existing models.Item
	if err := r.db.First(&existing, item.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("fetch item id=%d for update: %w", item.ID, err)
	}

	existing.Name = item.Name
	existing.Price = item.Price
	existing.Description = item.Description
	existing.SoldOut = item.SoldOut

	if err := r.db.Save(&existing).Error; err != nil {
		return nil, fmt.Errorf("update item id=%d: %w", item.ID, err)
	}

	return &existing, nil
}

func (r *itemRepository) Delete(itemID uint) error {
	result := r.db.Delete(&models.Item{}, itemID)
	if err := result.Error; err != nil {
		return fmt.Errorf("delete item id=%d: %w", itemID, err)
	}
	if result.RowsAffected == 0 {
		return ErrItemNotFound
	}
	return nil
}

// memoryItemRepository provides an in-memory implementation useful for tests.
type memoryItemRepository struct {
	mu     sync.RWMutex
	items  []models.Item
	nextID uint
}

func NewItemMemoryRepository(items []models.Item) IItemRepository {
	repo := &memoryItemRepository{
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

func (r *memoryItemRepository) FindAll() (*[]models.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]models.Item, len(r.items))
	copy(items, r.items)
	return &items, nil
}

func (r *memoryItemRepository) FindbyId(itemID uint) (*models.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.items {
		if r.items[i].ID == itemID {
			itemCopy := r.items[i]
			return &itemCopy, nil
		}
	}
	return nil, ErrItemNotFound
}

func (r *memoryItemRepository) Create(item models.Item) (*models.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if item.ID == 0 {
		item.ID = r.nextID
		r.nextID++
	} else {
		for _, existing := range r.items {
			if existing.ID == item.ID {
				return nil, fmt.Errorf("item id=%d already exists", item.ID)
			}
		}
		if item.ID >= r.nextID {
			r.nextID = item.ID + 1
		}
	}

	r.items = append(r.items, item)
	itemCopy := item
	return &itemCopy, nil
}

func (r *memoryItemRepository) Update(item models.Item) (*models.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.items {
		if r.items[i].ID == item.ID {
			r.items[i].Name = item.Name
			r.items[i].Price = item.Price
			r.items[i].Description = item.Description
			r.items[i].SoldOut = item.SoldOut

			itemCopy := r.items[i]
			return &itemCopy, nil
		}
	}
	return nil, ErrItemNotFound
}

func (r *memoryItemRepository) Delete(itemID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.items {
		if r.items[i].ID == itemID {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return ErrItemNotFound
}
