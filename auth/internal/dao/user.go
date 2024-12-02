package dao

import (
	"auth-ms/internal/dao/sql"
	u "auth-ms/internal/entity/user"
	"fmt"
)

// bool - is it already exist
func (s *Storage) SaveUser(usr u.SaveUser) (u.User, bool, error) {
	const fn = "dao.user_dao.SaveUser"

	preparedStatement, err := s.db.Prepare(sql.SQL_EXIST_CHECK_BY_USERID)

	if err != nil {
		return u.User{}, false, fmt.Errorf("%s: %w", fn, err)
	}

	var exist int

	res, err := preparedStatement.Query(usr.UserId)

	if err != nil {
		return u.User{}, false, fmt.Errorf("%s: %w", fn, err)
	}

	for res.Next() {
		err = res.Scan(&exist)

		if err != nil {
			return u.User{}, false, fmt.Errorf("%s: %w", fn, err)
		}
	}

	if exist != 1 {
		preparedStatement, err = s.db.Prepare(sql.SQL_INSERT_INTO_USER)

		if err != nil {
			return u.User{}, false, fmt.Errorf("%s: %w", fn, err)
		}

		execRes, err := preparedStatement.Exec(usr.UserId, usr.Name, usr.AvatarURL, usr.Sub)

		if err != nil {
			return u.User{}, false, fmt.Errorf("%s: %w", fn, err)
		}

		lastInsertId, err := execRes.LastInsertId()

		if err != nil {
			return u.User{}, false, fmt.Errorf("%s: %w", fn, err)
		}

		returnUser := u.User{
			Id:        int(lastInsertId),
			Name:      usr.Name,
			UserId:    usr.UserId,
			AvatarURL: usr.AvatarURL,
			Sub:       usr.Sub,
		}

		return returnUser, false, nil
	} else {
		preparedStatement, err := s.db.Prepare(sql.SQL_SELECT_USER_BY_USERID)

		if err != nil {
			return u.User{}, false, fmt.Errorf("%s: %w", fn, err)
		}

		res, err := preparedStatement.Query(usr.UserId)

		if err != nil {
			return u.User{}, false, fmt.Errorf("%s: %w", fn, err)
		}

		retUsr := u.User{}

		for res.Next() {
			var id int
			var uid string
			var name string
			var avatarURL string
			var sub u.Subscription

			err = res.Scan(&id, &uid, &name, &avatarURL, &sub)

			if err != nil {
				return u.User{}, false, fmt.Errorf("%s: %w", fn, err)
			}

			retUsr = u.User{
				Id:        id,
				UserId:    uid,
				Name:      name,
				AvatarURL: avatarURL,
				Sub:       sub,
			}
		}

		return retUsr, true, nil
	}
}

func (s *Storage) GetUser(userId string) (u.User, error) {
	const fn = "dao.user_dao.GetUser"

	preparedStatement, err := s.db.Prepare(sql.SQL_SELECT_USER_BY_USERID)

	if err != nil {
		return u.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	res, err := preparedStatement.Query(userId)

	if err != nil {
		return u.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	usr := u.User{}

	for res.Next() {
		var id int
		var uid string
		var name string
		var avatarURL string
		var sub u.Subscription

		err = res.Scan(&id, &uid, &name, &avatarURL, &sub)

		if err != nil {
			return u.User{}, fmt.Errorf("%s: %w", fn, err)
		}

		usr = u.User{
			Id:        id,
			UserId:    uid,
			Name:      name,
			AvatarURL: avatarURL,
			Sub:       sub,
		}
	}

	if usr == (u.User{}) {
		return u.User{}, fmt.Errorf("user doens`t exist")
	}

	return usr, nil
}

func (s *Storage) UpdateUser(user u.UpdateUser) (u.User, int64, error) {
	const fn = "dao.user.UpdateUser"

	preparedStatement, err := s.db.Prepare(sql.SQL_UPDATE_USER_BY_USERID)

	if err != nil {
		return u.User{}, 0, fmt.Errorf("%s: %w", fn, err)
	}

	execRes, err := preparedStatement.Exec(user.Name, user.AvatarURL, user.UserId)

	if err != nil {
		return u.User{}, 0, fmt.Errorf("%s: %w", fn, err)
	}

	rowsAffected, err := execRes.RowsAffected()

	if err != nil {
		return u.User{}, 0, fmt.Errorf("%s: %w", fn, err)
	}

	preparedStatement, err = s.db.Prepare(sql.SQL_SELECT_USER_BY_USERID)

	res, err := preparedStatement.Query(user.UserId)

	var respUser u.User
	for res.Next() {
		var id int
		var userid string
		var name string
		var avatarURL string
		var sub u.Subscription

		var err = res.Scan(&id, &userid, &name, &avatarURL, &sub)

		if err != nil {
			return u.User{}, 0, fmt.Errorf("%s: %w", fn, err)
		}
		respUser = u.User{
			Id:        id,
			Name:      name,
			UserId:    userid,
			AvatarURL: avatarURL,
			Sub:       sub,
		}
	}

	if err != nil {
		return u.User{}, 0, fmt.Errorf("%s: %w", fn, err)
	}

	return respUser, rowsAffected, nil
}

func (s *Storage) DeleteUser(user u.User) (int64, error) {
	const fn = "dao.user.DeleteUser"

	preparedStatement, err := s.db.Prepare(sql.SQL_SELECT_USER)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	var exist int

	res, err := preparedStatement.Query(user.Id, user.UserId, user.Name, user.AvatarURL, user.Sub)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	for res.Next() {
		err = res.Scan(&exist)

		if err != nil {
			return 0, fmt.Errorf("%s: %w", fn, err)
		}
	}

	if exist != 1 {
		return 0, fmt.Errorf("user doesn`t exist")
	}

	preparedStatement, err = s.db.Prepare(sql.SQL_DELETE_USER)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	execRes, err := preparedStatement.Exec(user.Id, user.UserId, user.Name, user.AvatarURL, user.Sub)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	rowsAffected, err := execRes.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return rowsAffected, nil
}
