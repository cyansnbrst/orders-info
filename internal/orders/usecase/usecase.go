package usecase

import (
	"log/slog"

	"cyansnbrst.com/order-info/config"
	"cyansnbrst.com/order-info/internal/models"
	"cyansnbrst.com/order-info/internal/orders"
	"cyansnbrst.com/order-info/internal/server/cache"
)

// Orders usecase
type ordersUC struct {
	cfg        *config.Config
	ordersRepo orders.Repository
	logger     *slog.Logger
	cache      *cache.Cache
}

// Orders usecase constructor
func NewOrdersUseCase(cfg *config.Config, ordersRepo orders.Repository, logger *slog.Logger, cache *cache.Cache) orders.UseCase {
	return &ordersUC{cfg: cfg, ordersRepo: ordersRepo, logger: logger, cache: cache}
}

// Get an order
func (u *ordersUC) Get(uid string) (*models.Order, error) {
	return u.ordersRepo.Get(uid)
}

// Save an order
func (u *ordersUC) Save(order *models.Order) error {
	return u.ordersRepo.Save(order)
}

// Get all orders
func (u *ordersUC) GetAll() ([]models.Order, error) {
	return u.ordersRepo.GetAll()
}
