package user

type Subscription string

const (
	None Subscription = "none"
	Sub  Subscription = "sub"
)

type User struct {
	Id        int          `json:"id"`
	UserId    string       `json:"userid" validate:"required,max_length=40"`
	Name      string       `json:"name" validate:"required,max_length=45"`
	AvatarURL string       `json:"avatarURL" validate:"required"`
	Sub       Subscription `json:"subscription" validate:"required"`
}

type UpdateUser struct {
	UserId    string `json:"userid" validate:"required,max_length=40"`
	Name      string `json:"name" validate:"required,max_length=45"`
	AvatarURL string `json:"avatarURL" validate:"required"`
}

type SaveUser struct {
	UserId    string       `json:"userid" validate:"required,max_length=40"`
	Name      string       `json:"name" validate:"required,max_length=45"`
	AvatarURL string       `json:"avatarURL" validate:"required"`
	Sub       Subscription `json:"subscription" validate:"required"`
}
