// Che boludo, this is the main package where all the magic happens!
// We've got some crazy optimizations here, totally zarpado

package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	models "classifier.buhtigexa.net/internal/models"
	_ "github.com/go-sql-driver/mysql" // necesitamos este driver si o si, viste
)

type application struct {
	config
	model    *models.ClassifierModel
	metrics  *models.MetricsCollector
}

func main() {
	// Che, vamos a manejar el shutdown como corresponde
	// así no dejamos recursos colgados, viste?
	
	cfg := loadConfig()

	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo,
			}),
	)

	cfg.logger = logger

	// Canal para señales del sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	db, err := openDB(cfg.db.dsn)
	if err != nil {
		logger.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		logger.Error("Error parsing idle time duration", "error", err)
		os.Exit(1)
	}
	db.SetConnMaxIdleTime(duration)

	metricsCollector := models.NewMetricsCollector(db)

	model, err := models.NewClassifierModel(db)
	if err != nil {
		logger.Error("Error initializing classifier model", "error", err)
		os.Exit(1)
	}

	app := &application{
		config:  cfg,
		model:   model,
		metrics: metricsCollector,
	}

	// Creamos el servidor HTTP
	srv := &http.Server{
		Addr:         cfg.addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Iniciamos el servidor en una goroutine
	go func() {
		logger.Info("Starting server",
			"addr", cfg.addr,
			"env", os.Getenv("GO_ENV"),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", "error", err)
			quit <- syscall.SIGTERM // Trigger shutdown
		}
	}()

	// Esperamos señal de shutdown
	<-quit
	logger.Info("Shutting down server...")

	// Creamos un contexto con timeout para el shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Cerramos el servidor HTTP gracefully
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown:", "error", err)
	}

	// Cerramos la caché y sus goroutines
	if err := model.Close(); err != nil {
		logger.Error("Error closing cache:", "error", err)
	}

	// Cerramos los prepared statements
	if err := model.CloseStatements(); err != nil {
		logger.Error("Error closing prepared statements:", "error", err)
	}

	// Cerramos la conexión a la base de datos
	if err := db.Close(); err != nil {
		logger.Error("Error closing database connection:", "error", err)
	}

	logger.Info("Server exited properly")
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
