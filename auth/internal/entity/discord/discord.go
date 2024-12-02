package discord

type Guild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GuildRoleMapping struct {
	GuildId   string   `yaml:"guildId"`
	GuildName string   `yaml:"guildName"`
	RoleId    []string `yaml:"roleId"`
}

type DiscordUser struct {
	Username  string `json:"username"`
	UserId    string `json:"userid"`
	AvatarURL string `json:"avatarURL"`
}
