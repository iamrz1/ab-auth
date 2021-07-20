package public

import (
	"github.com/go-chi/chi"
	"github.com/iamrz1/ab-auth/logger"
	"github.com/iamrz1/ab-auth/service"
)

type publicRouter struct {
	Services *service.Config
	Log      logger.Logger
}

func NewPublicRouter(svc *service.Config, logStruct logger.Logger) *publicRouter {
	return &publicRouter{
		Services: svc,
		Log:      logStruct,
	}
}

// Router returns a router
func (pr *publicRouter) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/customers", pr.customerRouter())
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
