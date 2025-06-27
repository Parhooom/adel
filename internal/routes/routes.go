package routes

import (
	"adel/internal/app"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(app.UserMiddleware.AuthenticateUser)

	r.Get("/health", app.HealthCheck)

	r.Route("/problems", func(r chi.Router) {
		r.Get("/{id}", app.ProblemHandler.HandleGetProblemByID)
		r.Get("/", app.ProblemHandler.HandleGetAllProblems)

		r.With(app.UserMiddleware.RequireAdminUser).Post("/", app.ProblemHandler.HandleCreateProblem)
		r.With(app.UserMiddleware.RequireAdminUser).Delete("/{id}", app.ProblemHandler.HandleDeleteProblemByID)
		r.With(app.UserMiddleware.RequireAdminUser).Put("/{id}", app.ProblemHandler.HandleUpdateProblemByID)
	})

	r.Route("/submissions", func(r chi.Router) {
		r.With(app.UserMiddleware.RequireUser).Post("/", app.SubmissionHandler.HandleCreateSubmission)
		r.With(app.UserMiddleware.RequireUser).Get("/{id}", app.SubmissionHandler.HandleGetSubmissionByID)
		r.With(app.UserMiddleware.RequireUser).Delete("/{id}", app.SubmissionHandler.HandleDeleteSubmissionByID)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/register", app.UserHandler.HandleRegisterNormalUser)
		r.Post("/login", app.TokenHandler.HandleLoginUser)
		r.Post("/register-admin", app.UserHandler.HandleRegisterAdminUser)

		r.With(app.UserMiddleware.RequireUser).Delete("/logout", app.TokenHandler.HandleDeleteToken)
	})

	return r
}
