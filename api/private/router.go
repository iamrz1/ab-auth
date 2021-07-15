package private

import (
	"github.com/go-chi/chi"
	"github.com/iamrz1/ab-auth/api/middleware"
	"github.com/iamrz1/ab-auth/logger"
	"github.com/iamrz1/ab-auth/service"
)

type privateRouter struct {
	Services *service.Config
	Log      logger.StructLogger
}

func NewPrivateRouter(svc *service.Config, logStruct logger.StructLogger) *privateRouter {
	return &privateRouter{
		Services: svc,
		Log:      logStruct,
	}
}

// Router returns a router
func (pr *privateRouter) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/customers", pr.customerRouter())
	return r
}

func (pr *privateRouter) customerRouter() *chi.Mux {
	r := chi.NewRouter()

	//r.Get("/", pr.ListGenerics)
	r.With(middleware.AuthenticatedOnly).Get("/profile", pr.GetCustomerProfile)
	r.With(middleware.AuthenticatedOnly).Patch("/profile", pr.UpdateCustomerProfile)
	r.With(middleware.AuthenticatedOnly).Get("/verify-token", pr.VerifyAccessToken)
	r.With(middleware.JWTTokenOnly).Get("/refresh-token", pr.RefreshToken)
	//r.Delete("/{slug}", pr.DeleteGeneric)
	//r.Delete("/purge/{slug}", pr.PurgeGeneric)

	return r
}
