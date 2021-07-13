package health

import (
	"github.com/go-chi/chi"
)

// Router return health router
func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", status)

	return r
}
