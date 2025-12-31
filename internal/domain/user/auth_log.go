package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/xiaohangshuhub/go-workit/pkg/ddd"
)

type UserAuthLog struct {
	ddd.AggregateRoot[uuid.UUID]           // 聚合根
	GrantType                    GrantType // 授权类型
	UserId                       uuid.UUID // 用户ID
	ClientID                     string    // 客户端ID
	Time                         time.Time // 授权时间
	Message                      *string   // 登录信息
	IPAddress                    *string   // 用户登录IP
	UserAgent                    *string   // 用户设备/浏览器信息
	Status                       bool      // 授权状态
}

// NewUserAuthLog 创建会话日志
func NewUserAuthLog(id, userId uuid.UUID, grantType GrantType, status bool) (*UserAuthLog, error) {

	if id == uuid.Nil {
		err := ErrSessionlogIdIsnil
		return nil, err
	}

	if userId == uuid.Nil {
		err := ErrSessionlogUserIdIsnil
		return nil, err
	}

	return &UserAuthLog{
		AggregateRoot: ddd.NewAggregateRoot(id),
		UserId:        userId,
		GrantType:     grantType,
		Status:        status,
	}, nil
}
