package user

// GrantType 授权类型
type GrantType int

const (
	AuthorizationCode   GrantType = iota // 授权码模式
	PasswordCredentials                  // 密码模式
	ClientCredentials                    // 客户端模式
	Implicit                             // 隐式模式
	Refreshing                           // 刷新模式
)
