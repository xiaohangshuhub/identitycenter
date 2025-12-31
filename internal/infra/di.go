package infra

import (
	"github.com/xiaohangshuhub/xiaohangshu/internal/infra/repoimpl"
	"go.uber.org/fx"
)

func DependencyInjection() []fx.Option {

	return []fx.Option{
		fx.Provide(repoimpl.NewUserRepository),
	}

}
