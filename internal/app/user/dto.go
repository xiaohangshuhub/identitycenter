package user

// UserInfoDTO 用户信息传输对象
type UserInfoDTO struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
}
