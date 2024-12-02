package dao

import (
	r "chat-ms/internal/entity/reactions"
	user_entity "chat-ms/internal/entity/user"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type ReactionManager interface {
	SaveUser(userId, name, avatarURL, guild string, contractAddress common.Address) error
	SetReaction(userId, guild string, value int, contractAddress common.Address) error
	GetReactions(userId, guild string, contractAddress common.Address) (r.Reactions, error)
	GetAllUserWithReactionsByAddressAndGuild(guild string, contractAddress common.Address) ([]user_entity.User, error)
}

func (s *Storage) SaveUser(userId, name, avatarURL, guild string, contractAddress common.Address) error {
	const fn = "dao.reactions.SaveUser"

	preparedStatement, err := s.db.Prepare(SQL_CHECK_REACTION)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	res, err := preparedStatement.Query(userId, contractAddress.Hex(), guild)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	var exist int

	for res.Next() {
		err = res.Scan(&exist)

		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	}

	if exist == 1 {
		return nil
	} else {
		preparedStatement, err = s.db.Prepare(SQL_INSERT_REACTION)

		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		_, err = preparedStatement.Exec(userId, name, avatarURL, 0, guild, contractAddress.Hex())
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	}

	return nil
}

func (s *Storage) GetReactions(userId, guild string, contractAddress common.Address) (r.Reactions, error) {
	const fn = "dao.reactions.GetReactions"

	reactions := r.Reactions{}

	preparedStatement, err := s.db.Prepare(SQL_GET_REACTION_BY_USERID_AND_CONTRACT_ADDRESS)

	if err != nil {
		return r.Reactions{}, fmt.Errorf("%s: %w", fn, err)
	}

	res, err := preparedStatement.Query(userId, contractAddress.Hex(), guild)

	if err != nil {
		return r.Reactions{}, fmt.Errorf("%s: %w", fn, err)
	}

	var reaction int
	for res.Next() {
		err = res.Scan(&reaction)

		if err != nil {
			return r.Reactions{}, fmt.Errorf("%s: %w", fn, err)
		}
	}

	reactions.Reaction = reaction

	preparedStatement, err = s.db.Prepare(SQL_GET_REACTIONS_BY_CONTRACT_ADDRESS)

	if err != nil {
		return r.Reactions{}, fmt.Errorf("%s: %w", fn, err)
	}

	res, err = preparedStatement.Query(contractAddress.Hex(), guild)

	if err != nil {
		return r.Reactions{}, fmt.Errorf("%s: %w", fn, err)
	}

	for res.Next() {
		err = res.Scan(&reaction)

		if err != nil {
			return r.Reactions{}, fmt.Errorf("%s: %w", fn, err)
		}

		if reaction <= -1 {
			reactions.DislikesAmount++
		} else if reaction >= 1 {
			reactions.LikesAmount++
		}
	}

	return reactions, nil
}

func (s *Storage) SetReaction(userId, guild string, value int, contractAddress common.Address) error {
	const fn = "dao.reactions.SetDislike"

	preparedStatement, err := s.db.Prepare(SQL_CHECK_REACTION)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	res, err := preparedStatement.Query(userId, contractAddress.Hex(), guild)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	var exist int

	for res.Next() {
		err = res.Scan(&exist)

		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	}

	if exist == 1 {
		preparedStatement, err = s.db.Prepare(SQL_UPDATE_REACTION)

		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		_, err = preparedStatement.Exec(value, userId, contractAddress.Hex(), guild)
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	} else {
		preparedStatement, err = s.db.Prepare(SQL_INSERT_REACTION)

		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		_, err = preparedStatement.Exec(userId, value, guild, contractAddress.Hex())
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	}

	return nil
}

func (s *Storage) GetAllUserWithReactionsByAddressAndGuild(guild string, contractAddress common.Address) ([]user_entity.User, error) {
	const fn = "dao.reactions.GetAllUserWithReactionsByAddressAndGuild"

	users := []user_entity.User{}

	preparedStatement, err := s.db.Prepare(SQL_SELECT_USERID_BY_CONTRACT_ADDRESS_AND_GUILD)

	if err != nil {
		return []user_entity.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	res, err := preparedStatement.Query(contractAddress.Hex(), guild)

	if err != nil {
		return []user_entity.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	var userid string
	var name string
	var avatarURL string
	var reaction int64
	for res.Next() {
		err = res.Scan(&userid, &name, &avatarURL, &reaction)

		if err != nil {
			return []user_entity.User{}, fmt.Errorf("%s: %w", fn, err)
		}

		users = append(users, user_entity.User{UserId: userid, Name: name, AvatarURL: avatarURL, Reaction: reaction})
	}

	return users, nil
}
