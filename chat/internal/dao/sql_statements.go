package dao

const (
	SQL_INSERT_MESSAGE = "INSERT INTO message (name, avatarURL, guild,  contractAddress, body,time,repliedMessageId, pinned) VALUES (?,?,?,?,?,?,?,?)"

	SQL_SELECT_MESSAGE_BY_CONTRACT_ADDRESS_AND_GUILD = "SELECT * FROM message WHERE  contractAddress = ? AND guild = ?"

	SQL_UPDATE_PINNED = "UPDATE message SET pinned = ? WHERE id = ?"

	SQL_GET_REACTION_BY_USERID_AND_CONTRACT_ADDRESS = "SELECT reaction FROM reaction WHERE userId = ? AND contractAddress = ? AND guild = ?"

	SQL_GET_REACTIONS_BY_CONTRACT_ADDRESS = "SELECT reaction FROM reaction WHERE contractAddress = ? AND guild = ?"

	SQL_CHECK_REACTION = "SELECT 1 FROM reaction WHERE userId = ? AND contractAddress = ? AND guild = ?"

	SQL_INSERT_REACTION = "INSERT INTO reaction (userId, name, avatarURL, reaction, guild, contractAddress) VALUE (?,?,?,?,?,?)"

	SQL_UPDATE_REACTION = "UPDATE reaction SET reaction = ? WHERE userId = ? AND contractAddress = ? AND guild = ?"

	SQL_SELECT_USERID_BY_CONTRACT_ADDRESS_AND_GUILD = "SELECT userId, name, avatarURL, reaction FROM reaction WHERE contractAddress = ? AND guild = ?"
)
