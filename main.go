package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/3-lines-studio/bifrost"
	webview "github.com/webview/webview_go"

	"github.com/3-lines-studio/datafrost/internal/adapter/database"
	adapterHttp "github.com/3-lines-studio/datafrost/internal/adapter/http"
	"github.com/3-lines-studio/datafrost/internal/adapter/repository"
	"github.com/3-lines-studio/datafrost/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var version = "dev"

//go:embed all:.bifrost
var bifrostFS embed.FS

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Datafrost %s\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) > 0 && flag.Args()[0] == "reset" {
		path := repository.DBPath()
		fmt.Printf("Resetting database at %s...\n", path)
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			log.Fatalf("Failed to remove database: %v", err)
		}
		configDB, err := repository.NewConfigDB()
		if err != nil {
			log.Fatalf("Failed to recreate database: %v", err)
		}
		configDB.Close()
		fmt.Println("Database reset successfully.")
		os.Exit(0)
	}

	configDB, err := repository.NewConfigDB()
	if err != nil {
		log.Fatalf("Failed to initialize config database: %v", err)
	}
	defer configDB.Close()

	sqlDB := configDB.DB()

	connectionRepo := repository.NewConnectionRepository(sqlDB)
	savedQueryRepo := repository.NewSavedQueryRepository(sqlDB)
	appStateRepo := repository.NewAppStateRepository(sqlDB)

	factory := database.NewFactory()
	adapterCache := database.NewAdapterCache()
	defer adapterCache.Close()

	connectionUsecase := usecase.NewConnectionUsecase(connectionRepo, factory, adapterCache)
	tableUsecase := usecase.NewTableUsecase(connectionRepo, adapterCache)
	queryUsecase := usecase.NewQueryUsecase(connectionRepo, adapterCache)
	savedQueryUsecase := usecase.NewSavedQueryUsecase(savedQueryRepo)
	appStateUsecase := usecase.NewAppStateUsecase(appStateRepo)
	adapterUsecase := usecase.NewAdapterUsecase(factory)

	connectionsHandler := adapterHttp.NewConnectionsHandler(connectionUsecase)
	tablesHandler := adapterHttp.NewTablesHandler(tableUsecase)
	queryHandler := adapterHttp.NewQueryHandler(queryUsecase)
	savedQueriesHandler := adapterHttp.NewSavedQueriesHandler(savedQueryUsecase)
	tabsHandler := adapterHttp.NewTabsHandler(appStateUsecase)
	themeHandler := adapterHttp.NewThemeHandler(appStateUsecase)
	layoutHandler := adapterHttp.NewLayoutHandler(appStateUsecase)
	adapterHandler := adapterHttp.NewAdapterHandler(adapterUsecase)

	apiRouter := chi.NewRouter()
	apiRouter.Use(middleware.Logger)
	apiRouter.Use(middleware.Recoverer)
	apiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

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
				r.Get("/tables/{name}/schema", tablesHandler.GetSchema)
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
		r.Get("/adapters", adapterHandler.List)
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

	setupMacEditMenu()

	w.Run()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
