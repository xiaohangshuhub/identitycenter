package user

import (
	"context"

	"github.com/xiaohangshuhub/xiaohangshu/internal/domain/user"
	"go.uber.org/zap"
)

// Login  登录请求结构体
type Login struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler  登录处理者
type LoginHandler struct {
	*zap.Logger
	repo user.UserRepository
}

// NewLoginHandler 创建一个查询处理者
func NewLoginHandler(repo user.UserRepository) *LoginHandler {
	return &LoginHandler{
		repo: repo,
	}
}

// Handle  处理登录请求,根据用户名和密码返回用户信息或错误信息。
func (h *LoginHandler) Handle(ctx context.Context, query *Login) (userDTO *UserInfoDTO, err error) {

	// 这里可以添加对用户名和密码的验证逻辑
	// 例如查询数据库验证用户凭据

	user, err := h.repo.GetUserInfoByPassword(ctx, query.Account, query.Password)

	if err != nil {
		return nil, err
	}

	// 假设用户名和密码验证通过，返回一个用户信息
	return &UserInfoDTO{
		UserID:    user.ID,
		Username:  query.Account,
		Nickname:  "Test User",
		AvatarURL: "http://example.com/avatar.jpg",
		Email:     "",
		Phone:     "",
	}, nil

}
