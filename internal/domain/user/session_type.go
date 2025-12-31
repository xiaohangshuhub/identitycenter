package user

type SessionType int

// 会话操作
const (
	Login  SessionType = iota //登录
	Logout                    //登出
)
