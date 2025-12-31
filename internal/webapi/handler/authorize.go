package handler

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/xiaohangshuhub/xiaohangshu/pkg/session"
	"go.uber.org/zap"
)

// Authorize godoc
// @Summary Authorize
// @Description 授权接口
// @Tags OAuth2
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /connect/authorize [get]
func Authorize(srv *server.Server, session *session.Session, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		w := c.Writer
		r := c.Request

		if v, _ := session.Get(r, "authorize_form"); v != nil {
			r.ParseForm()
			if r.Form.Get("client_id") == "" {
				r.Form = v.(url.Values)
			}
		}

		if err := session.Delete(w, r, "authorize_form"); err != nil {
			log.Error("delete request form error", zap.Error(err))
			return
		}

		if err := srv.HandleAuthorizeRequest(w, r); err != nil {
			c.JSON(500, ErrorResponse{Error: err.Error()})
			return
		}
	}
}
