package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/Gervva/avito_test_task/cmd"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	controller "github.com/Gervva/avito_test_task/internal/controller"
	bannerRepository "github.com/Gervva/avito_test_task/internal/repository"
	serviceBanner "github.com/Gervva/avito_test_task/internal/service/banner"
	storageCache "github.com/Gervva/avito_test_task/internal/storage/cache"
	storageDB "github.com/Gervva/avito_test_task/internal/storage/database"

	addBannerHandler "github.com/Gervva/avito_test_task/internal/handlers/add_banner"
	deleteBannerHandler "github.com/Gervva/avito_test_task/internal/handlers/delete_banner"
	getBannerHandler "github.com/Gervva/avito_test_task/internal/handlers/get_banner"
	getUserBannerHandler "github.com/Gervva/avito_test_task/internal/handlers/get_user_banner"
	updateBannerHandler "github.com/Gervva/avito_test_task/internal/handlers/update_banner"

	deleteByFeatureTagHandler "github.com/Gervva/avito_test_task/internal/handlers/delete_by_feature_tag"
	getAllVersionsHandler "github.com/Gervva/avito_test_task/internal/handlers/get_all_versions"
	getBannerVersionHandler "github.com/Gervva/avito_test_task/internal/handlers/get_banner_version"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	config, err := cmd.Load()
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.ErrorContext(context.Background(), "fail to connect to database", "error", err)
		return
	}
	// if err = db.Ping(); err != nil {
	// 	logger.ErrorContext(context.Background(), "fail to ping to database", "error", err)
	// 	return
	// }

	defer func() { _ = db.Close() }()

	database := storageDB.NewDatabase(db)
	databaseRepo := storageDB.NewRepository(database)

	cacheRepo := InitCache(*config)

	bannerRepo := bannerRepository.New(databaseRepo, cacheRepo)

	bannerService := serviceBanner.New(bannerRepo)

	handlerAddBanner := addBannerHandler.New(bannerService, logger)
	handlerDeleteBanner := deleteBannerHandler.New(bannerService, logger)
	handlerGetBanner := getBannerHandler.New(bannerService, logger)
	handlerGetUserBanner := getUserBannerHandler.New(bannerService, logger)
	handlerUpdateBanner := updateBannerHandler.New(bannerService, logger)

	handlerDeleteByFeatureTag := deleteByFeatureTagHandler.New(bannerService, logger)
	handlerGetAllVersions := getAllVersionsHandler.New(bannerService, logger)
	handlerGetBannerVersion := getBannerVersionHandler.New(bannerService, logger)

	r := mux.NewRouter()

	r.HandleFunc("/user_banner", controller.UserMW(handlerGetUserBanner)).Methods("GET")
	r.HandleFunc("/banner", controller.AdminMW(handlerGetBanner)).Methods("GET")
	r.HandleFunc("/banner", controller.AdminMW(handlerAddBanner)).Methods("POST")
	r.HandleFunc("/banner/{id}", controller.AdminMW(handlerUpdateBanner)).Methods("PATCH")
	r.HandleFunc("/banner/{id}", controller.AdminMW(handlerDeleteBanner)).Methods("DELETE")

	r.HandleFunc("/delete_by_feature_tag", controller.AdminMW(handlerDeleteByFeatureTag)).Methods("DELETE")
	r.HandleFunc("/get_all_versions", controller.AdminMW(handlerGetAllVersions)).Methods("GET")
	r.HandleFunc("/get_banner_version", controller.AdminMW(handlerGetBannerVersion)).Methods("GET")

	// staticHandler := http.StripPrefix(staticURIPrefix, http.FileServer(http.Dir(logCSVDirectory)))
	// mux.Handle(staticURIPrefix+"/", staticHandler)

	server := http.Server{
		Addr:    ":" + config.Microservice.Port,
		Handler: r,
	}

	go func() {
		logger.InfoContext(context.Background(), "service started", "port", config.Microservice.Port)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorContext(context.Background(), "error while starting server", "error", err)
		}
	}()

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	logger.InfoContext(context.Background(), "shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		logger.ErrorContext(context.Background(), "error while shutting down", "error", err)
	}
}

func InitCache(cfg cmd.Config) storageCache.Repository {
	db := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Cache.Host, cfg.Cache.Port),
	})

	cacheRepo := storageCache.NewRepository(db)

	return *cacheRepo
}
