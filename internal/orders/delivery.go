package orders

import "net/http"

// Orders handlers interface
type Handlers interface {
	Get() http.HandlerFunc
}
