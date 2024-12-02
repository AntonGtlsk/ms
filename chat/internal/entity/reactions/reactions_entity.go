package reactions


type Reactions struct{
	LikesAmount uint `json:"likes"`
	DislikesAmount uint `json:"dislikes"`
	Reaction int `json:"reaction"` // <=-1 - dislike 0 - empty >=1 - like
}