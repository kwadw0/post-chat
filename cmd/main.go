package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}

	cfg := config{
		addr: ":8000",
		db: dbConfig{
			dsn: os.Getenv("DB_DSN"),
		},
	}

	conn, err := pgx.Connect(context.Background(), cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	slog.Info("Connected to database successfully 🚀")

	m, err := migrate.New(
		"file://"+os.Getenv("MIGRATION_DIRECTORY"),
		cfg.db.dsn)
	if err != nil {
		slog.Error("migration failed: ", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("migration up failed: ", err)
	}
	slog.Info("Migrations applied successfully!")
	api := application{
		config: cfg,
	}
	//api.run(api.mount())

	//Default logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := api.run((api.mount())); err != nil {
		slog.Error("Application failed to start", "error", err)
		os.Exit(1)
	}
}
