package handlers

import (
	"chat-ms/internal/dao"
	h "chat-ms/internal/wss/hub"
	m "chat-ms/internal/wss/message"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"logging"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type URLParameters struct {
	AccessToken     string         `json:"access_token" validate:"required"`
	ContractAddress common.Address `json:"contract_address" validate:"required"`
	Guild           string         `json:"guild" validate:"required"`
	GuildId         string         `json:"guildid" validate:"required"`
	UserId          string         `json:"userid" validate:"required"`
	Name            string         `json:"name" validate:"required"`
	AvatarURL       string         `json:"avatarURL" validate:"required"`
}

func New(logging *logging.Logger, hub *h.Hub, messageManager dao.MessageManager, reactionManager dao.ReactionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		var urlParam URLParameters
		urlParam.AccessToken = r.URL.Query().Get("access_token")
		urlParam.Guild = r.URL.Query().Get("guild")
		urlParam.GuildId = r.URL.Query().Get("guildid")
		urlParam.ContractAddress = common.HexToAddress(r.URL.Query().Get("contractAddress"))
		urlParam.UserId = r.URL.Query().Get("userid")
		urlParam.Name = r.URL.Query().Get("name")
		urlParam.AvatarURL = r.URL.Query().Get("avatarURL")

		//isMember, err := is_guild_member.IsGuildMember(urlParam.Guild, urlParam.AccessToken, logging)
		//if !isMember || err != nil {
		//	fmt.Println(err)
		//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//	return
		//}

		reactions, err := reactionManager.GetReactions(urlParam.UserId, urlParam.Guild, urlParam.ContractAddress)

		if err != nil {
			http.Error(w, "Error getting reactions", http.StatusInternalServerError)
			logging.WithFields(logrus.Fields{
				"err": err,
			}).Error("Error getting reaction")
			return
		}

		messages, err := messageManager.GetGuildMessagesByAddress(urlParam.Guild, urlParam.ContractAddress)

		if err != nil {
			http.Error(w, "Error getting chat history", http.StatusInternalServerError)
			logging.WithFields(logrus.Fields{
				"err": err,
			}).Error("Error getting chat history")
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Error connection to wss")
		}

		client := &h.Client{UserId: urlParam.UserId,
			Hub: hub, Conn: conn, Send: make(chan m.Message, 256),
			ContractAddress: common.HexToAddress(r.URL.Query().Get("contractAddress")),
			GuildName:       urlParam.Guild,
			ChatHistory:     messages,
			Reactions:       reactions,
			Name:            urlParam.Name,
			AvatarURL:       urlParam.AvatarURL,
		}

		client.Hub.Register <- client

		go client.WritePump()
		go client.ReadPump(logging, messageManager, reactionManager)
	}
}
