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
	r.With(middleware.AuthenticatedOnly).Get("/profile", pr.getCustomerProfile)
	r.With(middleware.AuthenticatedOnly).Patch("/profile", pr.updateCustomerProfile)
	r.With(middleware.AuthenticatedOnly).Get("/verify-token", pr.verifyAccessToken)
	r.With(middleware.JWTTokenOnly).Get("/refresh-token", pr.refreshToken)
	r.With(middleware.AuthenticatedOnly).Put("/password", pr.updatePassword)

	r.Mount("/address", pr.addressRouter())

	//r.Delete("/{slug}", pr.DeleteGeneric)
	//r.Delete("/purge/{slug}", pr.PurgeGeneric)

	return r
}

func (pr *privateRouter) addressRouter() *chi.Mux {
	r := chi.NewRouter()

	// address apis, private
	r.With(middleware.AuthenticatedOnly).Post("/", pr.addNewAddressHandler)
	r.With(middleware.AuthenticatedOnly).Patch("/{id}", pr.updateAddressHandler)
	r.With(middleware.AuthenticatedOnly).Delete("/{id}", pr.removeAddressHandler) //empty body
	r.With(middleware.AuthenticatedOnly).Get("/all", pr.getAddressesHandler)
	r.With(middleware.AuthenticatedOnly).Get("/primary", pr.getPrimaryAddressHandler)
	r.With(middleware.AuthenticatedOnly).Post("/primary/{id}", pr.setPrimaryAddressHandler) //empty body

	//r.With(middleware.SecretOnly).Get("/secret/get-primary-address/{username}", getPrimaryAddressWithSecretHandler)

	return r
}
