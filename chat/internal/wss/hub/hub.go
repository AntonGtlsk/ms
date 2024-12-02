package hub

import (
	// api "chat-ms/internal/api/get_guilds_users"

	"chat-ms/internal/dao"
	"chat-ms/internal/entity/reactions"
	maxLengthValidator "chat-ms/internal/validator"
	m "chat-ms/internal/wss/message"
	"fmt"
	"logging"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 4096
)

type Client struct {
	Hub *Hub

	Conn *websocket.Conn

	Send chan m.Message

	ContractAddress common.Address

	GuildName string

	GuildId string

	UserId string

	Name string

	AvatarURL string

	Reactions reactions.Reactions

	ChatHistory []m.Message
}

type Hub struct {
	ClientsByContract map[common.Address]map[*Client]bool

	Broadcast chan m.Message

	Register chan *Client

	Unregister chan *Client
}

func (h *Hub) Run(reactionManager dao.ReactionManager, logger *logging.Logger) {
	for {
		select {
		case client := <-h.Register:
			if h.ClientsByContract[client.ContractAddress] == nil {
				h.ClientsByContract[client.ContractAddress] = make(map[*Client]bool)
			}
			h.ClientsByContract[client.ContractAddress][client] = true
			clients := h.ClientsByContract[client.ContractAddress]

			viewersCount := 0

			onlineClients := []Client{}

			for currentClient, _ := range h.ClientsByContract[client.ContractAddress] {
				if currentClient.GuildName == client.GuildName {
					onlineClients = append(onlineClients, *currentClient)
					viewersCount++
				}
			}

			err := reactionManager.SaveUser(client.UserId, client.Name, client.AvatarURL, client.GuildName, client.ContractAddress)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("Error save user")
			}

			reactionsClients, err := reactionManager.GetAllUserWithReactionsByAddressAndGuild(client.GuildName, client.ContractAddress)

			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("Error GetAllUserWithReactionsByAddressAndGuild user")
			}

			for i := range reactionsClients {
				for _, onlineClient := range onlineClients {
					if onlineClient.UserId == reactionsClients[i].UserId {
						reactionsClients[i].Online = true
					}
				}
			}

			client.Send <- m.Message{
				Event:           m.ChatStats,
				UserId:          client.UserId,
				GuildName:       client.GuildName,
				ContractAddress: client.ContractAddress,
				Users:           reactionsClients,
				Reactions:       client.Reactions,
				ChatHistory:     client.ChatHistory,
				Name:            client.GuildName,
				AvatarURL:       client.AvatarURL,
			}

			for existingClient := range clients {
				if existingClient.GuildName == client.GuildName {
					existingClient.Send <- m.Message{
						Event:           m.UserConnect,
						ContractAddress: client.ContractAddress,
						ViewersCount:    viewersCount,
						UserId:          client.UserId,
						Name:            client.GuildName,
						AvatarURL:       client.AvatarURL,
					}
				}
			}

		case client := <-h.Unregister:
			clients := h.ClientsByContract[client.ContractAddress]
			delete(h.ClientsByContract[client.ContractAddress], client)
			close(client.Send)

			viewersCount := 0

			for currentClient, _ := range h.ClientsByContract[client.ContractAddress] {
				if currentClient.GuildName == client.GuildName {
					viewersCount++
				}
			}

			for existingClient := range clients {
				if existingClient.GuildName == client.GuildName {
					existingClient.Send <- m.Message{
						Event:           m.UserDisconnect,
						ContractAddress: client.ContractAddress,
						ViewersCount:    viewersCount,
						UserId:          client.UserId,
						Name:            client.GuildName,
						AvatarURL:       client.AvatarURL,
					}
				}
			}

		case message := <-h.Broadcast:
			clients := h.ClientsByContract[message.ContractAddress]

			for client := range clients {
				if client.GuildName == message.GuildName {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.ClientsByContract[client.ContractAddress], client)
					}
				}
			}

		}
	}
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:         make(chan m.Message),
		Register:          make(chan *Client),
		Unregister:        make(chan *Client),
		ClientsByContract: make(map[common.Address]map[*Client]bool),
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteJSON(message)

			if err != nil {
				fmt.Println(err)
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump(logging *logging.Logger, messageManager dao.MessageManager, reactionManager dao.ReactionManager) {
	defer func() {
		c.Hub.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	saveMessage := &m.SaveMessage{}

	for {
		err := c.Conn.ReadJSON(saveMessage)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("unexpected close error:", err)
				break
			}

			break
		}

		v := validator.New()

		err = v.RegisterValidation("max_length", maxLengthValidator.MaxLengthValidator)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Warningf("Failed to validate data")
		}

		if err := v.Struct(saveMessage); err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Warningf("Failed to validate data")

			continue
		}

		switch saveMessage.Event {
		case m.NewMessage:

			message, err := messageManager.SaveMessage(saveMessage)

			if err != nil {
				fmt.Println(err)
			}

			c.Hub.Broadcast <- *message
		case m.PinMessage:

			message, err := messageManager.ChangePinnedMessage(saveMessage)
			if err != nil {
				fmt.Println(err)
			}

			c.Hub.Broadcast <- *message

		case m.Reaction:
			err = reactionManager.SetReaction(saveMessage.UserId, saveMessage.GuildName, saveMessage.Reactions.Reaction, saveMessage.ContractAddress)
			if err != nil {
				fmt.Println(err)
			}

			reactions, err := reactionManager.GetReactions(saveMessage.UserId, saveMessage.GuildName, saveMessage.ContractAddress)
			if err != nil {
				fmt.Println(err)
			}

			c.Hub.Broadcast <- m.Message{Event: m.ReactionsCount, Reactions: reactions, UserId: saveMessage.UserId, GuildName: saveMessage.GuildName, ContractAddress: saveMessage.ContractAddress}
		}

	}

}
