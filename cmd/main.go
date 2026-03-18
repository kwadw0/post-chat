// @title Post-Chat eCommerce API
// @version 1.0
// @description This is a sample server for a toy store.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @host localhost:8000
// @BasePath /

package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Setup the Brain (Logger) first!
	// We do this first so if the car won't start, we have a camera recording why.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. Read the Secret Map (.env)
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Could not find .env file, using system environment variables instead")
	}

	cfg := config{
		addr: ":8000",
		db: dbConfig{
			dsn: os.Getenv("DB_DSN"),
		},
		jwt_secret: []byte(os.Getenv("JWT_SECRET")),
		jwt_exp: func() time.Duration {
			val, _ := strconv.Atoi(os.Getenv("JWT_EXP"))
			if val == 0 {
				return 15 * time.Minute // Default to 15 mins
			}
			return time.Duration(val) * time.Minute
		}(),
	}

	// 3. Setup the Water Tank (Database Pool)
	// Instead of one straw, we create a box of straws (Pool) for the party!
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		slog.Error("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Check if we can actually talk to the database
	if err := pool.Ping(ctx); err != nil {
		slog.Error("Unable to ping database", "error", err)
		os.Exit(1)
	}

	slog.Info("Connected to database successfully 🚀")

	// 4. Update the House (Migrations)
	m, err := migrate.New(
		"file://"+os.Getenv("MIGRATION_DIRECTORY"),
		cfg.db.dsn)
	if err != nil {
		slog.Error("Migration setup failed", "error", err)
	} else {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			slog.Error("Migration up failed", "error", err)
		} else {
			slog.Info("Migrations applied successfully!")
		}
	}

	// 5. Start the Engine (Application)
	api := application{
		config: cfg,
		db:     pool, // Now we're giving it the real pool, not an empty box!
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Application failed to start", "error", err)
		os.Exit(1)
	}
}
