package message

import (
	"chat-ms/internal/entity/reactions"
	user_entity "chat-ms/internal/entity/user"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Event string

const (
	NewMessage     Event = "new_message"
	PinMessage     Event = "pin_message"
	UserConnect    Event = "user_connect"
	UserDisconnect Event = "user_disconnect"
	Reaction       Event = "reaction"
	ReactionsCount Event = "reactions_count"
	ChatStats      Event = "chat_stats"
)

type Message struct {
	Event        Event               `json:"event" validate:"required"`
	ViewersCount int                 `json:"viewers_count"` //when Event == UserConnect or UserDisconnect
	Reactions    reactions.Reactions `json:"reactions"`     //when Event == Reaction; <= -1 - dislike, 0 - nothing, >= 1 - like
	ChatHistory  []Message           `json:"chat_history"`
	Users        []user_entity.User  `json:"users"` //when Event == ChatStats
	Id           int64               `json:"id"`
	UserId       string              `json:"userid" validate:"required"`
	GuildName    string              `json:"guild_name" validate:"required"`
	// GuildId          string             `json:"guild_id" validate:"required"`
	Name             string         `json:"name" validate:"max_length=45"`
	AvatarURL        string         `json:"avatarURL" `
	ContractAddress  common.Address `json:"contract_address" validate:"required"`
	Body             string         `json:"body"`
	Time             time.Time      `json:"time"`
	RepliedMessageId int64          `json:"replied_message_id"`
	Pinned           *bool          `json:"pinned"`
}

type SaveMessage struct {
	Event        Event               `json:"event" validate:"required"`
	ViewersCount int                 `json:"viewers_count"` //when Event == UserConnect or UserDisconnect
	Reactions    reactions.Reactions `json:"reactions"`     //when Event == Reaction; <= -1 - dislike, 0 - nothing, >= 1 - like
	ChatHistory  []Message           `json:"chat_history"`
	Users        []user_entity.User  `json:"users"`         //when Event == UserStats
	Id           int64               `json:"id"`
	UserId       string              `json:"userid" validate:"required"`
	GuildName    string              `json:"guild_name" validate:"required"`
	// GuildId          string         `json:"guild_id" validate:"required"`
	Name             string         `json:"name" validate:"max_length=45"`
	AvatarURL        string         `json:"avatarURL"`
	ContractAddress  common.Address `json:"contract_address" validate:"required"`
	Body             string         `json:"body"`
	Time             time.Time      `json:"time"`
	RepliedMessageId int64          `json:"replied_message_id"`
	Pinned           *bool          `json:"pinned"`
}
