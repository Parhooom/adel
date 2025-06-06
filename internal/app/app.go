package app

import (
	"adel/internal/api"
	"adel/internal/service/postgres"
	"adel/internal/utils"
	"adel/migrations"

	"database/sql"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger            *log.Logger
	DB                *sql.DB
	ProblemHandler    *api.ProblemHandler
	SubmissionHandler *api.SubmissionHandler
}

func NewApplication() (*Application, error) {
	pgDB, err := postgres.Open()
	if err != nil {
		return nil, err
	}

	err = postgres.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// stores
	problemStore := postgres.NewPostgresProblemStore(pgDB)
	submissionStore := postgres.NewPostgresSubmissionStore(pgDB)

	// handlers
	problemHandler := api.NewProblemHandler(problemStore, logger)
	submissionHandler := api.NewSubmissionHandler(submissionStore, logger)

	app := &Application{
		Logger:            logger,
		DB:                pgDB,
		ProblemHandler:    problemHandler,
		SubmissionHandler: submissionHandler,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"status": "ok"})
}
