package server

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	ordersHttp "cyansnbrst.com/order-info/internal/orders/delivery/http"
	kafkaHandler "cyansnbrst.com/order-info/internal/orders/kafka"
	ordersRepository "cyansnbrst.com/order-info/internal/orders/repository"
	ordersUC "cyansnbrst.com/order-info/internal/orders/usecase"
	"cyansnbrst.com/order-info/internal/server/cache"
	kafkaConsumer "cyansnbrst.com/order-info/pkg/kafka"
)

// Register server handlers
func (s *Server) RegisterHandlers() http.Handler {
	router := httprouter.New()

	// Init repositories
	ordersRepo := ordersRepository.NewOrdersRepository(s.db)
	cache := cache.NewInMemoryCache()

	// Init cache
	if err := cache.Recover(ordersRepo); err != nil {
		s.logger.Error("error recovering cache from DB", "error", err)
	}
	cache.StartCleaner(5 * time.Minute)
	cache.PrintCache()

	// Init useCases
	ordersUC := ordersUC.NewOrdersUseCase(s.config, ordersRepo, s.logger, cache)

	// Init handlers
	ordersHandlers := ordersHttp.NewOrdersHandlers(s.config, ordersUC, s.logger, cache)

	// Register order routes
	ordersHttp.RegisterOrderRoutes(router, ordersHandlers)

	// Kafka consumer setup
	kafkaConsumer := kafkaConsumer.NewKafkaConsumer(s.config.Kafka.Brokers, s.config.Kafka.Topic, s.config.Kafka.Group, s.logger)
	kafkaHandler := kafkaHandler.NewKafkaMessageHandler(ordersUC, cache, s.logger)

	go func() {
		if err := kafkaConsumer.Consume(kafkaHandler.Handle); err != nil {
			s.logger.Error("error starting Kafka consumer", "error", err)
		}
	}()

	// Static files
	router.ServeFiles("/static/*filepath", http.Dir("./static"))

	return router
}
