package user_entity

type User struct{
	Name string `json:"name"`
	UserId string `json:"userid"`
	AvatarURL string `json:"avatarURL"`
	Reaction int64 `json:"reaction"` //<=-1 - dislike, 0 - no reaction, >=1 - like
	Online bool `json:"online"`
}
