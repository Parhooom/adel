package app

import (
	"adel/internal/api"
	"adel/internal/middleware"
	"adel/internal/service/judge"
	"adel/internal/service/postgres"
	"adel/internal/service/rabbitmq"
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
	UserHandler       *api.UserHandler
	TokenHandler      *api.TokenHandler
	UserMiddleware    *middleware.UserMiddleware
	JudgeService      *judge.JudgeService
	RabbitMQClient    *rabbitmq.RabbitMQClient
	RabbitMQWorkers   *rabbitmq.Worker
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

	dockerClient := judge.NewDockerClient(logger)
	judgeService := judge.NewJudge(dockerClient, logger)
	rabbitmqClient := rabbitmq.Open(logger)

	// stores
	problemStore := postgres.NewPostgresProblemStore(pgDB)
	submissionStore := postgres.NewPostgresSubmissionStore(pgDB)
	userStore := postgres.NewPostgresUserStore(pgDB)
	tokenStore := postgres.NewPostgresTokenStore(pgDB)

	// handlers
	problemHandler := api.NewProblemHandler(problemStore, logger)
	submissionHandler := api.NewSubmissionHandler(submissionStore, problemStore, rabbitmqClient, logger)
	userHandler := api.NewUserHandler(userStore, submissionStore, problemStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	// middlewares
	userMiddleware := &middleware.UserMiddleware{UserStore: userStore}

	rabbitmqWorkers := rabbitmq.NewRabbitMQWorker(rabbitmqClient, submissionStore, problemStore, 3, judgeService, logger)

	app := &Application{
		Logger:            logger,
		DB:                pgDB,
		ProblemHandler:    problemHandler,
		SubmissionHandler: submissionHandler,
		UserHandler:       userHandler,
		TokenHandler:      tokenHandler,
		UserMiddleware:    userMiddleware,
		JudgeService:      judgeService,
		RabbitMQClient:    rabbitmqClient,
		RabbitMQWorkers:   rabbitmqWorkers,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"status": "ok"})
}
