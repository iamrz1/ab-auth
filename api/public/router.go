package public

import (
	"github.com/go-chi/chi"
	"github.com/iamrz1/ab-auth/logger"
	"github.com/iamrz1/ab-auth/service"
)

type publicRouter struct {
	Services *service.Config
	Log      logger.StructLogger
}

func NewPublicRouter(svc *service.Config, logStruct logger.StructLogger) *publicRouter {
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

	r.Post("/signup", pr.Signup)
	r.Post("/verify-signup", pr.VerifySignUp)
	//r.Get("/", pr.ListGenerics)
	//r.Get("/{slug}", pr.GetGeneric)
	//r.Patch("/{slug}", pr.UpdateGeneric)
	//r.Delete("/{slug}", pr.DeleteGeneric)
	//r.Delete("/purge/{slug}", pr.PurgeGeneric)

	return r
}
