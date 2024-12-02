package main

import (
	"auth-ms/internal/config"
	"auth-ms/internal/dao"
	"auth-ms/internal/http-server/handlers/auth/discord/callback"
	"auth-ms/internal/http-server/handlers/auth/discord/login"
	"auth-ms/internal/http-server/handlers/auth/discord/logout"
	"auth-ms/internal/http-server/handlers/auth/discord/refresh"
	"auth-ms/internal/http-server/handlers/auth/discord/validate"
	"auth-ms/internal/http-server/handlers/sub/get"
	"auth-ms/internal/http-server/handlers/sub/update"
	getUserChatAccesss "auth-ms/internal/http-server/handlers/user_data/get_user_chat_access"
	getUserGuilds "auth-ms/internal/http-server/handlers/user_data/get_user_guilds"
	mw "auth-ms/internal/http-server/middleware/sub"
	"auth-ms/internal/service/auth"
	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"logging"
	"net/http"
)

const (
	all = "all"
)

func main() {
	cfg := config.MustLoad()

	logger := logging.New(cfg.BotToken, cfg.Logger.Folder, cfg.Logger.Filenames, cfg.Logger.WebhookUrls, cfg.Logger.UndefinedWebhook, all)

	jwtParser := auth.NewJwtParser([]byte(cfg.JWTKey))

	authMw := mw.NewAuth(*jwtParser)

	dg, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error creating discord client")
	}

	err = dg.Open()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error opening discord connection")
	}

	storage, err := dao.New(cfg.StoragePath())

	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Errorf("Failed to init storage")
	}

	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/discord/login", login.New(logger))
	router.Get("/discord/callback", callback.New(logger, storage, jwtParser))
	router.Get("/discord/refresh", refresh.New(logger, storage, jwtParser))
	router.Get("/discord/validate", validate.New(logger, jwtParser))
	router.Get("/discord/logout", logout.New(logger, storage))

	router.Get("/get-user-guilds", getUserGuilds.New(logger, storage, dg, jwtParser))

	router.Options("/get-user-guilds", getUserGuilds.NewOptions())

	router.Get("/get-chat-access", getUserChatAccesss.New(logger, dg, jwtParser))
	router.Options("/get-chat-access", getUserChatAccesss.NewOptions())

	router.Route("/sub", func(r chi.Router) {
		r.Use(authMw.New())
		r.Get("/get", get.New(logger, storage, jwtParser))
		r.Put("/update", update.New(logger, storage, *jwtParser))
	})

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	logger.Infof("starting server at %s", cfg.Address)

	if err := srv.ListenAndServe(); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Fatalf("failed to start server")
	}

	logger.Infof("server stopped")
}
