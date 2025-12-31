package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/xiaohangshuhub/go-workit/pkg/ddd"
)

// UserSessionLog 描述用户会话日志领域对象
type UserSessionLog struct {
	ddd.AggregateRoot[uuid.UUID]             // 聚合根
	SessionType                  SessionType // 会话类型
	UserId                       uuid.UUID   // 用户ID
	ClientID                     string      // 客户端ID
	Time                         time.Time   // 登录时间
	Message                      *string     // 登录信息
	IPAddress                    *string     // 用户登录IP
	UserAgent                    *string     // 用户设备/浏览器信息
	Status                       bool        // 登录状态
}

// NewUserSessionLog 创建用户会话日志
func NewUserSessionLog(id, userId uuid.UUID, sessionType SessionType, loginStatus bool, message, ipAddress, userAgent *string) (*UserSessionLog, error) {

	if id == uuid.Nil {
		err := ErrSessionlogIdIsnil
		return nil, err
	}

	if userId == uuid.Nil {
		err := ErrSessionlogUserIdIsnil
		return nil, err
	}

	return &UserSessionLog{
		AggregateRoot: ddd.NewAggregateRoot(id),
		UserId:        userId,
		Time:          time.Now(),
		Status:        loginStatus,
		Message:       message,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		SessionType:   sessionType,
	}, nil
}
