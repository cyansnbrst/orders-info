package orders

import "cyansnbrst.com/order-info/internal/models"

type UseCase interface {
	Get(uid string) (*models.Order, error)
	Save(order *models.Order) error
	GetAll() ([]models.Order, error)
}
