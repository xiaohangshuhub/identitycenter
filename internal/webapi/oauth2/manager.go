package oauth2

import (
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/xiaohangshuhub/xiaohangshu/configs"
)

func NewManager(cfg *configs.OAuth2, clientStore oauth2.ClientStore, accessGenerate oauth2.AccessGenerate, tokenStore oauth2.TokenStore) *manage.Manager {

	// Initialize manager
	mgr := manage.NewDefaultManager()

	// Token config
	mgr.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    time.Hour * time.Duration(cfg.Manager.AccessTokenExp),
		RefreshTokenExp:   time.Hour * 24 * 3,
		IsGenerateRefresh: true,
	})
	// Token store
	mgr.MapTokenStorage(tokenStore)

	// access token generator
	// 使用HS512算法生成JWT access token
	// 这里的Issuer和JWTSignedKey需要根据实际情况配置
	mgr.MapAccessGenerate(accessGenerate)

	mgr.MapClientStorage(clientStore)

	
	return mgr
}
