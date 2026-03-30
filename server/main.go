package main

import (
	"context"
	"log"
	"microservices/internal/db"
	"microservices/internal/handler"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	dbConn "microservices/pkg/db"
)



func main ()  {
	database := dbConn.NewDB()
	queries := db.New(database)

	logger := log.Default()
	mux := http.NewServeMux()

  taskHandler := handler.NewTaskHandler(queries)
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
		case http.MethodPost:
			taskHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed",http.StatusMethodNotAllowed)
		}
	})

	srv := &http.Server{
		Addr: ":9090",
		Handler: mux,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Println("Server Listening on port 9090")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			logger.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop  // block until signal is recieved
	logger.Println("Shutting Down Server........")

	// Time out for shutdown
	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Second)
	defer cancel()

	// Graceful shutdown
	if err := srv.Shutdown(ctx); err != nil{
		logger.Println("Forced shutdown:", err)
	}
	logger.Println("Server exited cleanly")

 }
