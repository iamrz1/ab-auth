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
	r.Get("/bd-area", pr.listBDArea)
	return r
}

func (pr *publicRouter) customerRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/signup", pr.signup)
	r.Post("/verify-signup", pr.verifySignUp)
	r.Post("/login", pr.login)
	r.Post("/forgot-password", pr.forgotPassword)
	r.Post("/set-password", pr.setPassword)
	//r.Get("/", pr.ListGenerics)
	//r.Delete("/{slug}", pr.DeleteGeneric)
	//r.Delete("/purge/{slug}", pr.PurgeGeneric)

	return r
}
