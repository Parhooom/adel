package main

import (
	"adel/internal/app"
	"adel/internal/routes"

	"flag"
	"fmt"
	"net/http"
	"time"

	_ "adel/docs"
)

// @title           Adel Online Judge API
// @version         1.0
// @description     This is a simple online judge
// @host      localhost:8080
// @BasePath  /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "backend server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	app.RabbitMQWorkers.Start()
	defer app.RabbitMQWorkers.Stop()

	r := routes.SetupRoutes(app)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("app started on %d\n", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
