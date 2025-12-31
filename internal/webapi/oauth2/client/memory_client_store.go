package client

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/xiaohangshuhub/xiaohangshu/configs"
)

// NewMemoryClientStore 创建内存客户端存储
func NewMemoryClientStore(cfg *configs.OAuth2) oauth2.ClientStore {
	// 内存client store
	clientStore := store.NewClientStore()

	// 设置客户端信息
	for _, v := range cfg.Clients {
		clientStore.Set(v.ID, &models.Client{
			ID:     v.ID,
			Secret: v.Secret,
			Domain: v.RedirectURIs[0], // 假设第一个重定向URI是主域名
		})
	}

	return clientStore
}
