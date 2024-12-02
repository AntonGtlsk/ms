package sql

const (
	SQL_UPDATE_SUBSCRIPTION_BY_USERID = "UPDATE user SET subscription = ? WHERE userid = ?"

	SQL_SELECT_SUBSCRIPTION_BY_USERID = "SELECT subscription FROM user WHERE userid = ?"
)
