package main

import (
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	models "classifier.buhtigexa.net/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	config
	model *models.ClassifierModel
}

func main() {
	var cfg config

	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelWarn,
			}),
	)

	dsn := "appuser:appusersecret@tcp(localhost:3306)/classifiersdb"
	db, err := openDB(dsn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	app := &application{
		config: config{
			logger: logger,
			addr:   ":4000",
		},
		model: &models.ClassifierModel{DB: db},
	}

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address i.e : -addr=:4000")
	flag.Parse()

	app.logger.Info("Starting server at :", slog.String("addr", cfg.addr))
	err = http.ListenAndServe(cfg.addr, app.routes())
	log.Fatal(err)
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
