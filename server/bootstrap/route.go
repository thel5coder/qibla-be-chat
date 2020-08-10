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
				r.Get("/login/{id}", userHandler.LoginHandler)
				r.Group(func(r chi.Router) {
					r.Use(mJwt.VerifyUserTokenCredential)
					r.Get("/token", userHandler.GetByTokenHandler)
					r.Get("/travelPackage", userHandler.GetTravelPackageHandler)
					r.Get("/jamaah/travelPackage/{id}", userHandler.GetJamaahHandler)
				})
			})

			fileHandler := api.FileHandler{Handler: handlerType}
			r.Route("/file", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(mJwt.VerifyUserTokenCredential)
					r.Post("/", fileHandler.UploadHandler)
				})
			})

			roomHandler := api.RoomHandler{Handler: handlerType}
			r.Route("/room", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(mJwt.VerifyUserTokenCredential)
					r.Get("/", roomHandler.FindAllByParticipantHandler)
					r.Get("/id/{id}", roomHandler.FindByIDHandler)
					r.Post("/", roomHandler.CreateHandler)
					r.Put("/id/{id}", roomHandler.UpdateHandler)
					r.Delete("/id/{id}", roomHandler.DeleteHandler)
				})
			})

			chatHandler := api.ChatHandler{Handler: handlerType}
			r.Route("/chat", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(mJwt.VerifyUserTokenCredential)
					r.Get("/", chatHandler.GetAllByRoomHandler)
					r.Post("/", chatHandler.CreateHandler)
					r.Delete("/id/{id}", chatHandler.DeleteHandler)
				})
			})

			participantHandler := api.ParticipantHandler{Handler: handlerType}
			r.Route("/participant", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(mJwt.VerifyUserTokenCredential)
					r.Post("/add", participantHandler.AddParticipantHandler)
					r.Post("/remove", participantHandler.RemoveParticipantHandler)
					r.Put("/leave/room/{id}", participantHandler.LeaveParticipantHandler)
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
	})
}
