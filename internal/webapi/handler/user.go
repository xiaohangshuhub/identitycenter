package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaohangshuhub/xiaohangshu/internal/app/user"
	"github.com/xiaohangshuhub/xiaohangshu/internal/webapi/response"
	"github.com/xiaohangshuhub/xiaohangshu/pkg/session"
	"go.uber.org/zap"
)

const (
	userIdTag = "user_id"
)

// Userinfo godoc
// @Summary Userinfo
// @Description 获取用户信息
// @Tags OAuth2
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /connect/userinfo [get]
func Userinfo(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := response.Success("hello userinfo")
		c.JSON(200, data)
	}
}

// Login godoc
// @Summary Login
// @Description 用户登录
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/user/login [post]
func Login(seesion *session.Session, userApp *user.UserApp, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		param := &user.Login{}

		if err := c.BindJSON(param); err != nil {
			logger.Error("bind json failed", zap.Error(err))
			c.JSON(400, ErrorResponse{Error: "login failed"})
		}

		data, err := userApp.LoginHandler.Handle(c, param)

		if err != nil {
			logger.Error("login failed", zap.Error(err))
			c.JSON(500, ErrorResponse{Error: "login failed"})
		}

		if err = seesion.Set(c.Writer, c.Request, userIdTag, data.UserID); err != nil {
			logger.Error("set user id to session failed", zap.Error(err))
			c.JSON(500, ErrorResponse{Error: "login failed"})
		}

		c.JSON(200, response.Success(data))
	}
}

// Logout godoc
// @Summary Logout
// @Description 用户登出
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/user/logout [post]
func Logout(seesion *session.Session, userApp *user.UserApp, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		userid, err := seesion.Get(c.Request, userIdTag)

		if err != nil {
			logger.Error("get user id from session failed", zap.Error(err))
			c.JSON(500, ErrorResponse{Error: "logout error"})
		}

		// logout log
		if userid != nil {
			userApp.LogoutHandler.Handle(c, &user.Logout{UserId: userid.(string)})
		}

		if err = seesion.Delete(c.Writer, c.Request, userIdTag); err != nil {
			logger.Error("delete user id from session failed", zap.Error(err))
			c.JSON(500, ErrorResponse{Error: "logout error"})
		}

		c.JSON(200, response.Success("logout success"))
	}
}
