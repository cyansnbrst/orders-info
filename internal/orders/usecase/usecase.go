package usecase

import (
	"log/slog"

	"cyansnbrst.com/order-info/config"
	"cyansnbrst.com/order-info/internal/models"
	"cyansnbrst.com/order-info/internal/orders"
	"cyansnbrst.com/order-info/internal/server/cache"
)

type ordersUC struct {
	cfg        *config.Config
	ordersRepo orders.Repository
	logger     *slog.Logger
	cache      *cache.Cache
}

func NewOrdersUseCase(cfg *config.Config, ordersRepo orders.Repository, logger *slog.Logger, cache *cache.Cache) orders.UseCase {
	return &ordersUC{cfg: cfg, ordersRepo: ordersRepo, logger: logger, cache: cache}
}

func (u *ordersUC) Get(uid string) (*models.Order, error) {
	return u.ordersRepo.Get(uid)
}

func (u *ordersUC) Save(order *models.Order) error {
	return u.ordersRepo.Save(order)
}

func (u *ordersUC) GetAll() ([]models.Order, error) {
	return u.ordersRepo.GetAll()
}
