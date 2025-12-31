package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"

	"github.com/xiaohangshuhub/xiaohangshu/internal/webapi/response"
	"go.uber.org/zap"
)

// Token godoc
// @Summary Token
// @Description 获取token
// @Tags OAuth2
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /connect/token [post]
func Token(srv *server.Server, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		err := srv.HandleTokenRequest(c.Writer, c.Request)
		if err != nil {
			log.Error("HandleTokenRequest error", zap.Error(err))
			c.JSON(500, ErrorResponse{Error: err.Error()})
		}
	}
}

// Revoke godoc
// @Summary Revoke
// @Description 撤销token
// @Tags OAuth2
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /connect/revoke [get]
func Revoke(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := response.Success("hello revoke")
		c.JSON(200, data)
	}
}

// Introspect godoc
// @Summary Introspect
// @Description 检查token有效性
// @Tags OAuth2
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /connect/introspect [get]
func Introspect(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := response.Success("hello introspect")
		c.JSON(200, data)
	}
}
