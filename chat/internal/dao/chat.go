package dao

import (
	m "chat-ms/internal/wss/message"
	"fmt"
	t "time"

	"github.com/ethereum/go-ethereum/common"
)

type MessageGetter interface {
	GetGuildMessagesByAddress(guild string, contractAddress common.Address) ([]m.Message, error)
}

func (s *Storage) GetGuildMessagesByAddress(guild string, contractAddress common.Address) ([]m.Message, error) {
	preparedStatement, err := s.db.Prepare(SQL_SELECT_MESSAGE_BY_CONTRACT_ADDRESS_AND_GUILD)

	if err != nil {
		fmt.Println(err)
	}

	var result []m.Message

	res, err := preparedStatement.Query(contractAddress.String(), guild)

	for res.Next() {

		var name string
		var avatarUrl string
		var id int64
		var guildFromDb string
		var contractAddressString string
		var body string
		var time t.Time
		var repliedMessageId int64
		var pinned bool

		err := res.Scan(&name, &avatarUrl, &contractAddressString, &guildFromDb, &body, &time, &repliedMessageId, &id, &pinned)

		if err != nil {
			fmt.Println(err)
		}

		result = append(result, m.Message{Id: id, Name: name, AvatarURL: avatarUrl, GuildName: guildFromDb, ContractAddress: common.HexToAddress(contractAddressString), Body: body, Time: time, RepliedMessageId: repliedMessageId, Pinned: &pinned})
	}

	return result, nil
}

type MessageManager interface {
	SaveMessage(messages *m.SaveMessage) (*m.Message, error)
	ChangePinnedMessage(message *m.SaveMessage) (*m.Message, error)
	GetGuildMessagesByAddress(guild string, contractAddress common.Address) ([]m.Message, error)
}

func (s *Storage) SaveMessage(message *m.SaveMessage) (*m.Message, error) {
	preparedStatement, err := s.db.Prepare(SQL_INSERT_MESSAGE)
	if err != nil {
		fmt.Println(err)
	}

	result := m.Message{}

	res, err := preparedStatement.Exec(message.Name, message.AvatarURL, message.GuildName, message.ContractAddress.String(), message.Body, message.Time, message.RepliedMessageId, message.Pinned)
	if err != nil {
		fmt.Println(err)
	}
	messageId, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}
	result = m.Message{Event: message.Event, Id: messageId, Name: message.Name, AvatarURL: message.AvatarURL, GuildName: message.GuildName, ContractAddress: message.ContractAddress, Body: message.Body, Time: message.Time, RepliedMessageId: message.RepliedMessageId, Pinned: message.Pinned}

	return &result, nil
}

func (s *Storage) ChangePinnedMessage(message *m.SaveMessage) (*m.Message, error) {
	preparedStatement, err := s.db.Prepare(SQL_UPDATE_PINNED)
	if err != nil {
		fmt.Println(err)
	}

	_, err = preparedStatement.Exec(*message.Pinned, message.Id)
	if err != nil {
		fmt.Println(err)
	}

	return &m.Message{Event: message.Event, Id: message.Id, Name: message.Name, AvatarURL: message.AvatarURL, GuildName: message.GuildName, ContractAddress: message.ContractAddress,
		Body: message.Body, Time: message.Time, RepliedMessageId: message.RepliedMessageId, Pinned: message.Pinned}, nil
}
