package user

import (
	"context"

	"go.uber.org/zap"
)

// Logout  登出请求结构体
type Logout struct {
	UserId string `json:"user_id" binding:"required"`
}

// LogoutHandler  登出处理者
type LogoutHandler struct {
	log *zap.Logger
}

// NewLogoutHandler 创建一个查询处理者
func NewLogoutHandler(log *zap.Logger) *LogoutHandler {
	return &LogoutHandler{
		log: log,
	}
}

// Handle  处理登出请求
func (h *LogoutHandler) Handle(ctx context.Context, query *Logout) {

}
