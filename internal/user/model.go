package user

type UserInfor struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	FullName string `json:"full_name"`
	Roles    []Role `json:"roles"`
	Avartar  string `json:"avatar"`
}

type Role struct {
	RoleID   string `json:"role_id" bson:"role_id"`
	RoleName string `json:"role_name" bson:"role_name"`
}
