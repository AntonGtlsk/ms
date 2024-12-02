package dao

import (
	"auth-ms/internal/dao/sql"
)

func (s *Storage) IsUserExist(userId string) (bool, error) {

	preparedStatement, err := s.db.Prepare(sql.SQL_SELECT_USER_BY_USERID)

	if err != nil {
		return false, err
	}

	res, err := preparedStatement.Query(userId)

	if err != nil {
		return false, err
		//logging of pushing error?
	}

	if res.Next() {
		return true, nil
	}

	return false, nil
}
