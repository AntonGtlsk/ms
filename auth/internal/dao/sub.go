package dao

import (
	"auth-ms/internal/dao/sql"
	u "auth-ms/internal/entity/user"
	"fmt"
)

func (s *Storage) UpdateSubscription(userid string, subscription u.Subscription) (int64, error) {
	const fn = "dao.sub.SaveSubscriptionToUser"

	preparedStatement, err := s.db.Prepare(sql.SQL_UPDATE_SUBSCRIPTION_BY_USERID)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	res, err := preparedStatement.Exec(subscription, userid)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return rowsAffected, nil
}

func (s *Storage) GetSubscription(userid string) (u.Subscription, error) {
	const fn = "dao.sub.GetSubscription"

	preparedStatement, err := s.db.Prepare(sql.SQL_SELECT_SUBSCRIPTION_BY_USERID)

	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	res, err := preparedStatement.Query(userid)

	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	var sub u.Subscription

	for res.Next() {
		err = res.Scan(&sub)

		if err != nil {
			return "", fmt.Errorf("%s: %w", fn, err)
		}
	}

	return sub, nil
}
