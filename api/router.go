package api

import (
	"github.com/go-chi/chi"
	"github.com/iamrz1/ab-auth/api/public"
	"github.com/iamrz1/ab-auth/logger"
	"github.com/iamrz1/ab-auth/service"
)

func V1Router(svc *service.Config, logger logger.StructLogger) *chi.Mux {
	r := chi.NewRouter()
	publicRouter := public.NewPublicRouter(svc, logger)
	r.Mount("/public", publicRouter.Router())

	return r
}
