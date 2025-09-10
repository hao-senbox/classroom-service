package user

type UserInfor struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Avartar  Avatar `json:"avatar"`
}

type Avatar struct {
	ImageID  uint64 `json:"image_id"`
	ImageKey string `json:"image_key"`
	ImageUrl string `json:"image_url"`
	Index    int    `json:"index"`
	IsMain   bool   `json:"is_main"`
}
