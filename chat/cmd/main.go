package main

import (
	c "chat-ms/internal/config"
	"chat-ms/internal/dao"
	chatHistory "chat-ms/internal/http/handlers/chat_history"
	getReactions "chat-ms/internal/http/handlers/reactions/get_reactions"
	setReaction "chat-ms/internal/http/handlers/reactions/set_reaction"

	// setLike "chat-ms/internal/http/handlers/reactions/set_like"
	wssHandlers "chat-ms/internal/wss/handlers"
	h "chat-ms/internal/wss/hub"
	"logging"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	all = "all"
)

func main() {

	config := c.MustLoad()

	logger := logging.New(config.BotToken, config.Logger.Folder, config.Logger.Filenames, config.Logger.WebhookUrls, config.Logger.UndefinedWebhook, all)

	storage := dao.MustConnect(config.StoragePath(), logger)

	hub := h.NewHub()
	go hub.Run(storage, logger)

	r := mux.NewRouter()
	r.HandleFunc("/socket", wssHandlers.New(logger, hub, storage, storage))

	r.HandleFunc("/get-history", chatHistory.New(logger, storage)).Methods("GET")
	r.HandleFunc("/get-history", chatHistory.NewOptions()).Methods("OPTIONS")

	r.HandleFunc("/get-reactions", getReactions.New(logger, storage)).Methods("GET")
	r.HandleFunc("/get-reactions", getReactions.NewOptions()).Methods("OPTIONS")

	r.HandleFunc("/set-reaction", setReaction.New(logger, storage)).Methods("POST")
	r.HandleFunc("/set-reaction", setReaction.NewOptions()).Methods("OPTIONS")
	logger.Infof("starting server at %s", config.WSSServer.Address)
	err := http.ListenAndServe(config.WSSServer.Address, r)
	if err != nil {
		logger.Fatalf("Error starting server:", err)
	}
}
