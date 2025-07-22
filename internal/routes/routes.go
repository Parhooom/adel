package routes

import (
	"adel/internal/app"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func corsMiddleware() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(corsMiddleware())
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.UserMiddleware.AuthenticateUser)

	r.Get("/health", app.HealthCheck)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

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
		r.With(app.UserMiddleware.RequireUser).Get("/user", app.SubmissionHandler.HandleGetSubmissionsByUserID)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/register", app.UserHandler.HandleRegisterNormalUser)
		r.Post("/login", app.TokenHandler.HandleLoginUser)
		r.Post("/register-admin", app.UserHandler.HandleRegisterAdminUser)

		r.With(app.UserMiddleware.RequireUser).Get("/me", app.UserHandler.HandleGetCurrentUser)
		r.With(app.UserMiddleware.RequireUser).Get("/stats", app.UserHandler.HandleGetUserStats)
		r.With(app.UserMiddleware.RequireUser).Delete("/logout", app.TokenHandler.HandleDeleteToken)

		r.With(app.UserMiddleware.RequireAdminUser).Get("/", app.UserHandler.HandleGetAllUsers)
		r.With(app.UserMiddleware.RequireAdminUser).Delete("/{id}", app.UserHandler.HandleDeleteUser)
	})

	r.Route("/admin", func(r chi.Router) {
		r.With(app.UserMiddleware.RequireAdminUser).Get("/stats", app.UserHandler.HandleGetAdminStats)
		r.With(app.UserMiddleware.RequireAdminUser).Get("/problems", app.ProblemHandler.HandleGetAllProblemsForAdmin)
	})

	return r
}
