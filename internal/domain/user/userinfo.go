package user

type Status int

const (
	Normal  Status = iota // 正常
	Disable               // 禁用
)

type UserInfo struct {
	ID        string // 用户ID
	Loginname string // 用户名
	Nickname  string // 昵称
	Avatar    string // 头像地址
	Status    Status // 是否启用
}
