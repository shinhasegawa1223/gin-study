package dto

type CreateItemInput struct {
	Name        string `json:"name" binding:"required,min=2"`
	Price       int    `json:"price" binding:"required,min=1,max=99999999"`
	Description string `json:"description" binding:"required"`
}
