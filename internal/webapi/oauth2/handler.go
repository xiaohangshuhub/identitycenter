package oauth2

import (
	"context"
	"net/http"
	"strings"

	"github.com/xiaohangshuhub/xiaohangshu/configs"
	"github.com/xiaohangshuhub/xiaohangshu/internal/domain/user"
	"github.com/xiaohangshuhub/xiaohangshu/pkg/session"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/server"
	"go.uber.org/zap"
)

type OAuth2Handlers struct {
	session *session.Session
	*zap.Logger
	cfg  *configs.OAuth2
	repo user.UserRepository
}

func NewOAuth2Handlers(session *session.Session, cfg *configs.OAuth2, repo user.UserRepository, logger *zap.Logger) *OAuth2Handlers {
	return &OAuth2Handlers{
		session: session,
		Logger:  logger,
		cfg:     cfg,
		repo:    repo,
	}
}

// passwordAuthorizationHandler 密码授权
func (h *OAuth2Handlers) passwordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {

	// 参数校验
	if username == "" || password == "" {
		h.Error("passwordAuthorizationHandler Error: username or password is empty")
		return "", errors.ErrInvalidRequest
	}

	// 这里可以添加对用户名和密码的验证逻辑
	user, err := h.repo.GetUserInfoByPassword(ctx, username, password)

	if err != nil {
		h.Error("passwordAuthorizationHandler Error: username or password is invalid ", zap.String("username", username))
		return "", errors.ErrInvalidGrant
	}

	// 假设用户名和密码验证通过，返回一个用户ID
	return user.ID, nil
}

// userAuthorizeHandler 用户授权
func (h *OAuth2Handlers) userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {

	// 读取会话中的用户ID
	v, _ := h.session.Get(r, "user_id")

	// 如果会话中没有用户ID，重定向到登录页面
	if v == nil {

		// 如果请求的表单数据为空，解析表单,这一步是为了确保在登录页面可以获取到表单数据
		if r.Form == nil {
			r.ParseForm()
		}

		// 将请求的表单数据存入会话,这样在登录页面可以获取到用户之前的请求数据
		h.session.Set(w, r, "authorize_form", r.Form)

		// 登录页面最终会把userId写进session(user_id)
		w.Header().Set("Location", h.cfg.LoginURL)

		w.WriteHeader(http.StatusFound)

		return
	}

	// 如果会话中有用户ID，直接返回
	userID = v.(string)

	// 不记住用户
	//h.session.Delete(w, r, "user_id")

	return
}

// authorizeScopeHandler 授权域处理
func (h *OAuth2Handlers) authorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	if r.Form == nil {
		r.ParseForm()
	}
	// 从上下文读取 client_id 和 scope
	clinet_id := r.Form.Get("client_id")
	scope = r.Form.Get("scope")

	if clinet_id == "" {
		h.Error("authorizeScopeHandler Error: client_id is empty")
		return "", errors.ErrInvalidRequest
	}

	if scope == "" {
		h.Error("authorizeScopeHandler Error: scope is empty")
		return "", errors.ErrInvalidRequest
	}

	// 通过client_id获取client
	client, err := h.cfg.GetClient(clinet_id)

	if err != nil {
		h.Error("authorizeScopeHandler Error: client_id is invalid  ", zap.String("client_id", clinet_id))
		return "", errors.ErrInvalidClient
	}

	// 校验scope是否合法
	for _, v := range strings.Split(scope, " ") {
		// 过滤空格
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}

		// 验证scope是否在client注册的scope中
		if !client.ContainsScope(v) {
			h.Error("authorizeScopeHandler Error: scope is invalid ", zap.String("client_id", clinet_id), zap.String("scope", v))
			// 请求域不存在属于客户端未授权
			return "", errors.ErrInvalidScope
		}
	}

	return
}

// clientAuthorizedHandler 客户端授权
func (h *OAuth2Handlers) clientAuthorizedHandler(clientID string, grant oauth2.GrantType) (allowed bool, err error) {

	// 这里可以添加对client_id和grant_type的验证逻辑

	client, err := h.cfg.GetClient(clientID)

	if err != nil {
		h.Error("clientAuthorizedHandler Error:  client_id is invalid ", zap.String("client_id", clientID))
		return false, errors.ErrInvalidClient
	}

	if !client.ContainsGrantType(string(grant)) {
		h.Error("clientAuthorizedHandler Error: grant_type is invalid ", zap.String("client_id", clientID), zap.String("grant_type", grant.String()))
		// 不包含属于客户端为授权
		return false, errors.ErrUnauthorizedClient
	}

	return true, nil
}

// internalErrorHandler 内部错误处理
func (h *OAuth2Handlers) internalErrorHandler(err error) (re *errors.Response) {
	h.Error("internalErrorHandler Error:", zap.Error(err))
	re = errors.NewResponse(err, errors.StatusCodes[err])
	re.Description = errors.Descriptions[err]
	return
}

// responseErrorHandler 响应错误处理
func (h *OAuth2Handlers) responseErrorHandler(re *errors.Response) {
	h.Error("responseErrorHandler Error:", zap.String("error", re.Error.Error()), zap.String("description", re.Description))
}

// preRedirectErrorHandler 预重定向错误处理
func (h *OAuth2Handlers) preRedirectErrorHandler(w http.ResponseWriter, req *server.AuthorizeRequest, err error) error {

	if req != nil {

		h.Error("preRedirectErrorHandler Error:", zap.String("redirect uri", req.RedirectURI))
	}

	return err
}

// extensionFieldsHandler 扩展字段处理
func (h *OAuth2Handlers) extensionFieldsHandler(ti oauth2.TokenInfo) (fieldsValue map[string]interface{}) {

	fieldsValue = make(map[string]interface{})

	fieldsValue["Issuer"] = h.cfg.Issuer

	return fieldsValue
}
