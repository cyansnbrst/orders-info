package http

import (
	"net/http"

	"cyansnbrst.com/order-info/internal/orders"
	"github.com/julienschmidt/httprouter"
)

func RegisterOrderRoutes(router *httprouter.Router, h orders.Handlers) {
	router.HandlerFunc(http.MethodGet, "/order/:uid", h.Get())
}
