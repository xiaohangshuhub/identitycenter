package webapi

import (
	"github.com/xiaohangshuhub/xiaohangshu/configs"
	"github.com/xiaohangshuhub/xiaohangshu/internal/app"
	"github.com/xiaohangshuhub/xiaohangshu/internal/infra"
	"github.com/xiaohangshuhub/xiaohangshu/internal/webapi/oauth2"
	"github.com/xiaohangshuhub/xiaohangshu/internal/webapi/oauth2/client"
	"github.com/xiaohangshuhub/xiaohangshu/internal/webapi/oauth2/token"
	"github.com/xiaohangshuhub/xiaohangshu/pkg/http"
	"github.com/xiaohangshuhub/xiaohangshu/pkg/session"
	"go.uber.org/fx"
)

func DependencyInjection() []fx.Option {
	di := []fx.Option{
		fx.Provide(oauth2.NewManager),
		fx.Provide(oauth2.NewOAuth2Service),
		fx.Provide(oauth2.NewOAuth2Handlers),
		fx.Provide(http.NewUserHTTPClient),
		fx.Provide(client.NewMemoryClientStore),
		fx.Provide(token.NewCustomJWTAccessGenerate),
		fx.Provide(token.NewMemotyTokenStore),
		fx.Provide(session.NewSession),
	}

	di = append(di, configs.DependencyInjection()...)
	di = append(di, app.DependencyInjection()...)
	di = append(di, infra.DependencyInjection()...)

	return di
}
