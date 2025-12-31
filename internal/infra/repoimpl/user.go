package repoimpl

import (
	"context"

	"github.com/xiaohangshuhub/xiaohangshu/internal/domain/user"
	"github.com/xiaohangshuhub/xiaohangshu/pkg/http"
)

type UserRepositoryImpl struct {
	userClient *http.UserClient
}

// GetUserInfoByUsername 根据用户名获取用户信息
func (impl *UserRepositoryImpl) GetUserInfoByPassword(ctx context.Context, username, password string) (*user.UserInfo, error) {

	//请求用户服务
	//_, err = impl.userClient.GetUserInfoByPwd(ctx, username, password)

	if username != "admin" || password != "admin" {

		err := user.ErrInvalidUserOrPassword

		return nil, err
	}

	// 假设用户名和密码验证通过，返回一个用户信息
	return &user.UserInfo{
		ID:        "1",
		Loginname: username,
		Nickname:  "Test User",
		Avatar:    "http://example.com/avatar.jpg",
	}, nil
}

func NewUserRepository(client *http.UserClient) user.UserRepository {
	return &UserRepositoryImpl{
		userClient: client,
	}
}
