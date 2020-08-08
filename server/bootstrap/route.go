package bootstrap

import (
	"qibla-backend-chat/pkg/logruslogger"
	api "qibla-backend-chat/server/handler"
	"qibla-backend-chat/server/middleware"

	chimiddleware "github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RegisterRoutes ...
func (boot *Bootup) RegisterRoutes() {
	handlerType := api.Handler{
		DB:         boot.DB,
		EnvConfig:  boot.EnvConfig,
		Validate:   boot.Validator,
		Translator: boot.Translator,
		ContractUC: &boot.ContractUC,
		Jwe:        boot.Jwe,
		Jwt:        boot.Jwt,
	}
	mJwt := middleware.VerifyMiddlewareInit{
		ContractUC: &boot.ContractUC,
	}

	boot.R.Route("/v1", func(r chi.Router) {
		// Define a limit rate to 1000 requests per IP per request.
		rate, _ := limiter.NewRateFromFormatted("1000-S")
		store, _ := sredis.NewStoreWithOptions(boot.ContractUC.Redis, limiter.StoreOptions{
			Prefix:   "limiter_global",
			MaxRetry: 3,
		})
		rateMiddleware := stdlib.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))
		r.Use(rateMiddleware.Handler)

		// Logging setup
		r.Use(chimiddleware.RequestID)
		r.Use(logruslogger.NewStructuredLogger(boot.EnvConfig["LOG_FILE_PATH"], boot.EnvConfig["LOG_DEFAULT"]))
		r.Use(chimiddleware.Recoverer)

		// API
		r.Route("/api", func(r chi.Router) {
			userHandler := api.UserHandler{Handler: handlerType}
			r.Route("/user", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(mJwt.VerifyJwtTokenCredential)
					r.Get("/token", userHandler.GetByTokenHandler)
				})
			})

			odooHandler := api.OdooHandler{Handler: handlerType}
			r.Route("/odoo", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Get("/field/{id}", odooHandler.GetFieldHandler)
					r.Get("/travelPackage", odooHandler.FindAllTravelPackageHandler)
					r.Get("/travelPackage/{id}", odooHandler.FindByIDTravelPackageHandler)
					r.Get("/partner/{id}", odooHandler.FindByIDPartnerHandler)
					r.Get("/guide/{id}", odooHandler.FindByIDGuideHandler)
				})
			})
		})

		// API ADMIN
		r.Route("/api-admin", func(r chi.Router) {
			adminHandler := api.AdminHandler{Handler: handlerType}
			r.Route("/admin", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Post("/login", adminHandler.LoginHandler)
				})
				r.Group(func(r chi.Router) {
					r.Use(mJwt.VerifySuperadminTokenCredential)
					r.Get("/", adminHandler.GetAllHandler)
					r.Get("/id/{id}", adminHandler.GetByIDHandler)
					r.Get("/code/{code}", adminHandler.GetByCodeHandler)
					r.Post("/", adminHandler.CreateHandler)
					r.Put("/id/{id}", adminHandler.UpdateHandler)
					r.Delete("/id/{id}", adminHandler.DeleteHandler)
				})
			})
		})
	})
}
