package orders

import "net/http"

type Handlers interface {
	Get() http.HandlerFunc
}
