package main

import (
	repo "kwadw0/gocommerce/internal/adapters/postgres/sqlc"
	"kwadw0/gocommerce/internal/products"
	"kwadw0/gocommerce/users"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "kwadw0/gocommerce/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (app *application) run(h http.Handler) error {
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	slog.Info("Application has started listening", "addr", app.config.addr)
	return srv.ListenAndServe()
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi Test"))
	})

	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.GetAllProducts)
	r.Post("/products", productHandler.AddProduct)

	userService := users.NewUserService(repo.New(app.db))
	userHandler := users.NewHandler(userService)
	r.Post("/signup", userHandler.CreateUser)

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusSeeOther)
	})
	r.Get("/docs/*", httpSwagger.WrapHandler)
	return r
}

type application struct {
	config config

	db *pgxpool.Pool
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
