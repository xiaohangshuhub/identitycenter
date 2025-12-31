package webapi

import (
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/xiaohangshuhub/xiaohangshu/configs"
	"github.com/xiaohangshuhub/xiaohangshu/internal/app/user"
	"github.com/xiaohangshuhub/xiaohangshu/internal/webapi/handler"
	"github.com/xiaohangshuhub/xiaohangshu/pkg/session"
	"go.uber.org/zap"
)

// login 登录跳转授权页面
func login(r *gin.Engine) {
	r.GET("/login", func(ctx *gin.Context) {

		// 加载登录页面
		ctx.File("../web/auth/dist/index.html")
	})
}

// 授权端口:V1
func authApiV1EndPoint(r *gin.Engine, srv *server.Server, cfg *configs.OAuth2, session *session.Session, logger *zap.Logger) {

	connect := r.Group("connect")
	{
		connect.GET("authorize", handler.Authorize(srv, session, logger))
		connect.POST("token", handler.Token(srv, logger))
		connect.GET("userinfo", handler.Userinfo(logger))
		connect.GET("introspect", handler.Introspect(logger))
		connect.GET("revoke", handler.Revoke(logger))
	}

	wellknownGroup := r.Group(".well-known")
	{
		wellknownGroup.GET("openid-configuration/jwks", handler.Jwks(cfg, logger))
		wellknownGroup.GET("openid-configuration", handler.OpenidConfiguration(cfg, logger))
	}

}

// 用户端口:V1
func userApiV1EndPoint(r *gin.Engine, userApp *user.UserApp, session *session.Session, logger *zap.Logger) {

	user := r.Group("user")
	{
		user.POST("login", handler.Login(session, userApp, logger))
		user.POST("logout", handler.Logout(session, userApp, logger))
	}
}

var EndPointList = []any{
	login,
	authApiV1EndPoint,
	userApiV1EndPoint,
}
