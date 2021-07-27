package private

import (
	"github.com/go-chi/chi"
	"github.com/iamrz1/ab-auth/service"
	rLog "github.com/iamrz1/rest-log"
)

type privateRouter struct {
	Services *service.Config
	Log      rLog.Logger
}

func NewPrivateRouter(svc *service.Config, rLogger rLog.Logger) *privateRouter {
	return &privateRouter{
		Services: svc,
		Log:      rLogger,
	}
}

// Router returns a router
func (pr *privateRouter) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/customers", pr.customerRouter())
	return r
}
