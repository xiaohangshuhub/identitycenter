package token

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/store"
	"go.uber.org/zap"
)

// NewMemotyTokenStore 创建一个内存TokenStore
//
// 参数:
//
//	logger *zap.Logger 日志对象
//
// 返回值:
//
//	oauth2.TokenStore 内存TokenStore
func NewMemotyTokenStore(logger *zap.Logger) oauth2.TokenStore {

	tokenStore, err := store.NewMemoryTokenStore()

	if err != nil {
		logger.Panic("Failed to create memory token store", zap.Error(err))
		panic(err)
	}

	return tokenStore
}
