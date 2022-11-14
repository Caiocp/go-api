package database

import "github.com/caiocp/go-api/internal/entities"

type UserInterface interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
}

type ProductInterface interface {
	Create(product *entities.Product) error
	FindAll(page, limit int, sort string) ([]entities.Product, error)
	FindByID(id string) (*entities.Product, error)
	Update(product *entities.Product) error
	Delete(id string) error
}
