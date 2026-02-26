package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/3-lines-studio/bifrost"
	webview "github.com/webview/webview_go"

	"datafrost/internal/db"
	"datafrost/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

//go:embed all:.bifrost
var bifrostFS embed.FS

func main() {
	configDB, err := db.NewConfigDB()
	if err != nil {
		log.Fatalf("Failed to initialize config database: %v", err)
	}
	defer configDB.Close()

	connectionStore := db.NewConnectionStore(configDB.DB())

	apiRouter := chi.NewRouter()
	apiRouter.Use(middleware.Logger)
	apiRouter.Use(middleware.Recoverer)
	apiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	connectionsHandler := handlers.NewConnectionsHandler(connectionStore)
	tablesHandler := handlers.NewTablesHandler(connectionStore)
	queryHandler := handlers.NewQueryHandler(connectionStore)
	themeHandler := handlers.NewThemeHandler(configDB.DB())
	layoutHandler := handlers.NewLayoutHandler(configDB.DB())
	tabsHandler := handlers.NewTabsHandler(configDB.DB())
	savedQueriesStore := db.NewSavedQueriesStore(configDB.DB())
	savedQueriesHandler := handlers.NewSavedQueriesHandler(savedQueriesStore)

	apiRouter.Route("/api", func(r chi.Router) {
		r.Route("/connections", func(r chi.Router) {
			r.Get("/", connectionsHandler.List)
			r.Post("/", connectionsHandler.Create)
			r.Route("/{id}", func(r chi.Router) {
				r.Delete("/", connectionsHandler.Delete)
				r.Put("/", connectionsHandler.Update)
				r.Post("/select", connectionsHandler.SetLastConnected)
				r.Post("/test", connectionsHandler.TestExisting)
				r.Get("/tables", tablesHandler.List)
				r.Get("/tables/{name}", tablesHandler.GetData)
				r.Post("/query", queryHandler.Execute)
				r.Get("/tabs", tabsHandler.Get)
				r.Post("/tabs", tabsHandler.Save)
				r.Route("/queries", func(r chi.Router) {
					r.Get("/", savedQueriesHandler.List)
					r.Post("/", savedQueriesHandler.Create)
					r.Route("/{queryId}", func(r chi.Router) {
						r.Put("/", savedQueriesHandler.Update)
						r.Delete("/", savedQueriesHandler.Delete)
					})
				})
			})
			r.Post("/test", connectionsHandler.Test)
		})
		r.Get("/theme", themeHandler.Get)
		r.Post("/theme", themeHandler.Update)
		r.Get("/layouts/{key}", layoutHandler.Get)
		r.Post("/layouts/{key}", layoutHandler.Save)
	})

	app := bifrost.New(
		bifrostFS,
		bifrost.Page("/", "./web/app.tsx", bifrost.WithClient()),
	)
	defer app.Stop()

	server := &http.Server{
		Handler: app.Wrap(apiRouter),
		Addr:    "127.0.0.1:0",
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	localURL := fmt.Sprintf("http://%s", listener.Addr().String())

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	debug := os.Getenv("BIFROST_DEV") == "1"

	w := webview.New(debug)
	defer w.Destroy()

	w.SetTitle("Datafrost")
	w.SetSize(1200, 800, webview.HintNone)
	w.Navigate(localURL)

	w.Run()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
