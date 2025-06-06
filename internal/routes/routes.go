package routes

import (
	"adel/internal/app"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)

	r.Route("/problems", func(r chi.Router) {
		r.Get("/{id}", app.ProblemHandler.HandleGetProblemByID)
		r.Post("/", app.ProblemHandler.HandleCreateProblem)
		r.Delete("/{id}", app.ProblemHandler.HandleDeleteProblemByID)
		r.Put("/{id}", app.ProblemHandler.HandleUpdateProblemByID)
	})

	r.Route("/submissions", func(r chi.Router) {
		r.Post("/", app.SubmissionHandler.HandleCreateSubmission)
		r.Get("/{id}", app.SubmissionHandler.HandleGetSubmissionByID)
		r.Delete("/{id}", app.SubmissionHandler.HandleDeleteSubmissionByID)
	})

	return r
}
