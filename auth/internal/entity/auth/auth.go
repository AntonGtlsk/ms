package auth

type RefreshSession struct {
	Id                  int
	Userid              string
	RefreshToken        string
	DiscordRefreshToken string
	ExpiresIn           int64
	CreatedAt           int64
}
