package http

import (
	"log/slog"
	"net/http"

	"cyansnbrst.com/order-info/config"
	"cyansnbrst.com/order-info/internal/orders"
	"cyansnbrst.com/order-info/internal/server/cache"
	"cyansnbrst.com/order-info/pkg/utils"
)

type ordersHandlers struct {
	cfg           *config.Config
	ordersUseCase orders.UseCase
	logger        *slog.Logger
	cache         *cache.Cache
}

func NewOrdersHandlers(cfg *config.Config, ordersUseCase orders.UseCase, logger *slog.Logger, cache *cache.Cache) orders.Handlers {
	return &ordersHandlers{
		cfg:           cfg,
		ordersUseCase: ordersUseCase,
		logger:        logger,
		cache:         cache,
	}
}

func (h *ordersHandlers) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := utils.ReadUIDParam(r)
		if uid == "" {
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
				"error": "missing uid param",
			}, nil)
			return
		}
		if cachedOrder, found := h.cache.Get(uid); found {
			utils.WriteJSON(w, http.StatusOK, utils.Envelope{
				"order": cachedOrder,
			}, nil)
			return
		}

		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{
			"error": "order not found",
		}, nil)
	}
}
