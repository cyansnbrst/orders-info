package kafka

import (
	"encoding/json"
	"log/slog"
	"time"

	"cyansnbrst.com/order-info/internal/models"
	"cyansnbrst.com/order-info/internal/orders"
	"cyansnbrst.com/order-info/internal/server/cache"
)

// Kafka message handler
type KafkaMessageHandler struct {
	ordersUC orders.UseCase
	cache    *cache.Cache
	logger   *slog.Logger
}

// Kafka message handler constructor
func NewKafkaMessageHandler(ordersUC orders.UseCase, cache *cache.Cache, logger *slog.Logger) *KafkaMessageHandler {
	return &KafkaMessageHandler{
		ordersUC: ordersUC,
		cache:    cache,
		logger:   logger,
	}
}

// Kafka order message handler
func (h *KafkaMessageHandler) Handle(msg []byte) error {
	var order models.Order

	err := json.Unmarshal(msg, &order)
	if err != nil {
		h.logger.Error("failed to unmarshal message", "error", err)
		return err
	}

	err = h.ordersUC.Save(&order)
	if err != nil {
		h.logger.Error("failed to save order to DB", "error", err)
		return err
	}

	h.cache.Set(order.OrderUID, order, 5*time.Minute)

	return nil
}
