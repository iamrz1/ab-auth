package public

import (
	"github.com/go-chi/chi"
	"github.com/iamrz1/ab-auth/service"
	rLog "github.com/iamrz1/rest-log"
)

type publicRouter struct {
	Services *service.Config
	Log      rLog.Logger
}

func NewPublicRouter(svc *service.Config, rLogger rLog.Logger) *publicRouter {
	return &publicRouter{
		Services: svc,
		Log:      rLogger,
	}
}

// Router returns a router
func (pr *publicRouter) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/customers", pr.customerRouter())
	r.Mount("/merchants", pr.merchantRouter())
	r.Get("/bd-area", pr.listBDArea)
	return r
}
