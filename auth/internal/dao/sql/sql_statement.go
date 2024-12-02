package sql

const (
	SQL_INSERT_INTO_USER = "INSERT INTO user (userId, name, avatarURL, subscription) VALUES (?,?,?,?)"

	SQL_SELECT_USER = "SELECT 1 FROM user WHERE id = ? AND userId = ? AND name = ? AND avatarURL = ? AND access = ?"

	SQL_EXIST_CHECK_BY_USERID = "SELECT 1 FROM user WHERE userId = ?"

	SQL_SELECT_USER_BY_USERID = "SELECT * FROM user WHERE userId = ?"

	SQL_DELETE_USER = "DELETE FROM user WHERE id = ? AND userId = ? AND name = ?  AND avatarURL = ? AND access = ?"

	SQL_UPDATE_USER_BY_USERID = "UPDATE user SET  name = ?, avatarURL = ? WHERE userId = ?"
)
