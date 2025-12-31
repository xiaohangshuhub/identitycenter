package app

import (
	"github.com/xiaohangshuhub/xiaohangshu/internal/app/user"
	"go.uber.org/fx"
)

func DependencyInjection() []fx.Option {

	return []fx.Option{
		fx.Provide(user.NewUserApp),
		fx.Provide(user.NewLoginHandler),
		fx.Provide(user.NewLogoutHandler),
	}

}
