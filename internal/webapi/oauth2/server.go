package oauth2

import (
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
)

func NewOAuth2Service(mgr *manage.Manager, handler *OAuth2Handlers) *server.Server {

	// Create server
	srv := server.NewServer(server.NewConfig(), mgr)
	srv.SetInternalErrorHandler(handler.internalErrorHandler)                 // 内部错误处理
	srv.SetResponseErrorHandler(handler.responseErrorHandler)                 // 响应错误处理
	srv.SetAuthorizeScopeHandler(handler.authorizeScopeHandler)               // 作用域处理
	srv.SetUserAuthorizationHandler(handler.userAuthorizeHandler)             // 用户授权处理
	srv.SetPasswordAuthorizationHandler(handler.passwordAuthorizationHandler) //密码授权处理
	srv.SetClientAuthorizedHandler(handler.clientAuthorizedHandler)           // 客户端授权
	srv.SetPreRedirectErrorHandler(handler.preRedirectErrorHandler)           // 重定向前的错误处理
	//srv.SetExtensionFieldsHandler(handler.extensionFieldsHandler)             // 扩展字段
	return srv
}
