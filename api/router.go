package api

import (
	"github.com/go-chi/chi"
	"github.com/iamrz1/ab-auth/api/private"
	"github.com/iamrz1/ab-auth/api/public"
	"github.com/iamrz1/ab-auth/service"
	rLog "github.com/iamrz1/rest-log"
)

func V1Router(svc *service.Config, logger rLog.Logger) *chi.Mux {
	r := chi.NewRouter()
	publicRouter := public.NewPublicRouter(svc, logger)
	privateRouter := private.NewPrivateRouter(svc, logger)
	r.Mount("/public", publicRouter.Router())
	r.Mount("/private", privateRouter.Router())

	return r
}
