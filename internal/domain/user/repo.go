package user

import (
	"context"
)

type (
	UserRepository interface {

		// GetUserInfoByPassword 根据用户名和密码获取用户信息
		GetUserInfoByPassword(ctx context.Context, username, password string) (*UserInfo, error)
	}
)
