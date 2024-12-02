package dao

import (
	ea "auth-ms/internal/entity/auth"
	"fmt"
)

func (s *Storage) SetRefreshSession(userId string, refreshSession ea.RefreshSession) error {
	const fn = "dao.refresh.SetRefreshSession"

	preparedStatement, err := s.db.Prepare("INSERT INTO refreshsessions (userId, refreshToken, discordRefreshToken, expiresIn, createdAt) VALUES(?,?,?,?,?)")

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = preparedStatement.Exec(userId, refreshSession.RefreshToken, refreshSession.DiscordRefreshToken, refreshSession.ExpiresIn, refreshSession.CreatedAt)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) GetRefreshSession(token string) (ea.RefreshSession, error) {
	const fn = "dao.refresh.GetRefreshSession"

	preparedStatement, err := s.db.Prepare("SELECT * FROM refreshsessions WHERE refreshToken = ?")

	if err != nil {
		return ea.RefreshSession{}, fmt.Errorf("%s: %w", fn, err)
	}

	var session ea.RefreshSession

	res, err := preparedStatement.Query(token)

	if err != nil {
		return ea.RefreshSession{}, fmt.Errorf("%s: %w", fn, err)
	}

	for res.Next() {
		var id int
		var userid string
		var refreshToken string
		var discordRefreshToken string
		var expiresIn int64
		var createdAt int64
		res.Scan(&id, &userid, &refreshToken, &discordRefreshToken, &expiresIn, &createdAt)
		session = ea.RefreshSession{
			Id:                  id,
			Userid:              userid,
			RefreshToken:        refreshToken,
			DiscordRefreshToken: discordRefreshToken,
			ExpiresIn:           expiresIn,
			CreatedAt:           createdAt,
		}
	}

	if session == (ea.RefreshSession{}) {
		return ea.RefreshSession{}, fmt.Errorf("Session doesn`t exist")
	}

	return session, nil
}

func (s *Storage) DeleteRefreshSession(token string) (int64, error) {
	const fn = "dao.refresh.DeleteSession"

	preparedStatement, err := s.db.Prepare("DELETE FROM refreshsessions WHERE refreshToken = ?")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	res, err := preparedStatement.Exec(token)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return rowsAffected, nil
}

func (s *Storage) GetDiscordRefreshToken(userid string) (string, error) {
	const fn = "dao.refresh.GetDiscordRefreshToken"

	preparedStatement, err := s.db.Prepare("SELECT discordRefreshToken FROM refreshsessions WHERE userId = ?")

	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	res, err := preparedStatement.Query(userid)

	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	var discordRefreshToken string

	for res.Next() {
		err = res.Scan(&discordRefreshToken)

		if err != nil {
			return "", fmt.Errorf("%s: %w", fn, err)
		}
	}

	if discordRefreshToken == "" {
		if err != nil {
			return "", fmt.Errorf("Discord refresh token is empty, error: %s", err)
		}
	}

	return discordRefreshToken, nil
}

func (s *Storage) UpdateDiscordRefreshToken(oldToken, newToken string) (int64, error) {
	const fn = "dao.refresh.UpdateDiscordRefreshToken"

	preparedStatements, err := s.db.Prepare("UPDATE refreshsessions SET discordRefreshToken = ? WHERE discordRefreshToken = ?")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	res, err := preparedStatements.Exec(newToken, oldToken)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return rowsAffected, nil

}
