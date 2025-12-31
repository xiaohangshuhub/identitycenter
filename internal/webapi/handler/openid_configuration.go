package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/xiaohangshuhub/xiaohangshu/configs"
	"github.com/xiaohangshuhub/xiaohangshu/internal/webapi/response"
	"go.uber.org/zap"
)

type (

	// OpenIDConfiguration OpenID Connect 配置结构体
	OpenIDConfiguration struct {
		Issuer                                 string   `json:"issuer"`                                               // 发行者标识符 (iss)
		AuthorizationEndpoint                  string   `json:"authorization_endpoint"`                               // 授权端点 URL
		TokenEndpoint                          string   `json:"token_endpoint"`                                       // 令牌端点 URL
		UserinfoEndpoint                       string   `json:"userinfo_endpoint,omitempty"`                          // 用户信息端点 URL
		JwksURI                                string   `json:"jwks_uri"`                                             // 公钥集合 (JWKS) URL
		RegistrationEndpoint                   string   `json:"registration_endpoint,omitempty"`                      // 客户端注册端点（可选）
		DeviceAuthorizationEndpoint            string   `json:"device_authorization_endpoint,omitempty"`              // 设备授权端点（可选）
		IntrospectionEndpoint                  string   `json:"introspection_endpoint,omitempty"`                     // Token Introspection 端点（可选）
		RevocationEndpoint                     string   `json:"revocation_endpoint,omitempty"`                        // Token 撤销端点（可选）
		ResponseTypesSupported                 []string `json:"response_types_supported"`                             // 支持的响应类型
		SubjectTypesSupported                  []string `json:"subject_types_supported"`                              // 支持的 Subject 类型
		IDTokenSigningAlgValuesSupported       []string `json:"id_token_signing_alg_values_supported"`                // ID Token 签名算法
		UserinfoSigningAlgValuesSupported      []string `json:"userinfo_signing_alg_values_supported,omitempty"`      // 用户信息签名算法
		AuthorizationSigningAlgValuesSupported []string `json:"authorization_signing_alg_values_supported,omitempty"` // 授权响应签名算法
		ScopesSupported                        []string `json:"scopes_supported"`                                     // 支持的 Scope
		GrantTypesSupported                    []string `json:"grant_types_supported,omitempty"`                      // 支持的授权方式
		ClaimsSupported                        []string `json:"claims_supported,omitempty"`                           // 支持的 Claims
		CodeChallengeMethodsSupported          []string `json:"code_challenge_methods_supported,omitempty"`           // 支持的 PKCE 方法
		TokenEndpointAuthMethodsSupported      []string `json:"token_endpoint_auth_methods_supported,omitempty"`      // Token 端点认证方式
		DeviceCodeChallengeMethodsSupported    []string `json:"device_code_challenge_methods_supported,omitempty"`    // 设备端点支持的 PKCE 方法
		ClaimsParameterSupported               bool     `json:"claims_parameter_supported,omitempty"`                 // 是否支持 claims 参数
		RequirePushedAuthorizationRequests     bool     `json:"require_pushed_authorization_requests,omitempty"`      // 是否强制使用 PAR
		FrontchannelLogoutSupported            bool     `json:"frontchannel_logout_supported,omitempty"`              // 是否支持前端通道登出
		FrontchannelLogoutSessionSupported     bool     `json:"frontchannel_logout_session_supported,omitempty"`      // 是否支持前端登出时会话管理
		BackchannelLogoutSupported             bool     `json:"backchannel_logout_supported,omitempty"`               // 是否支持后端通道登出
		BackchannelLogoutSessionSupported      bool     `json:"backchannel_logout_session_supported,omitempty"`       // 是否支持后端登出时会话管理
	}
)

// OpenidConfiguration godoc
// @Summary OpenID 配置发现端点
// @Description 提供 OpenID Connect 发现文档
// @Tags WellKnown
// @Accept json
// @Produce json
// @Success 200 {object} OpenIDConfiguration
// @Router /.well-known/openid-configuration [get]
func OpenidConfiguration(cfg *configs.OAuth2, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		issuer := cfg.Issuer

		if issuer == "" {
			log.Error("OpenID Connect issuer not configured")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server configuration error"})
			return
		}

		// 配置中获取JWT签名方法
		signingMethod := "RS256" // 默认值
		if cfg.Manager != nil && cfg.Manager.SigningMethod != "" {
			signingMethod = cfg.Manager.SigningMethod
		}

		data := OpenIDConfiguration{
			Issuer:                            issuer,
			AuthorizationEndpoint:             issuer + "/connect/authorize",
			TokenEndpoint:                     issuer + "/connect/token",
			UserinfoEndpoint:                  issuer + "/connect/userinfo",
			JwksURI:                           issuer + "/.well-known/openid-configuration/jwks",
			ResponseTypesSupported:            []string{"code", "token", "id_token"},
			SubjectTypesSupported:             []string{"public"},
			IDTokenSigningAlgValuesSupported:  []string{signingMethod},
			ScopesSupported:                   []string{"openid", "profile", "email", "offline_access"},
			ClaimsSupported:                   []string{"sub", "name", "email", "email_verified"},
			GrantTypesSupported:               []string{"authorization_code", "refresh_token", "password", "client_credentials"},
			TokenEndpointAuthMethodsSupported: []string{"client_secret_basic", "client_secret_post"},
			CodeChallengeMethodsSupported:     []string{"S256", "plain"},
		}

		c.JSON(http.StatusOK, data)
	}
}

// Jwks godoc
// @Summary JWKS 端点
// @Description 提供签名密钥的JWK集合
// @Tags WellKnown
// @Accept json
// @Produce json
// @Success 200 {object} jwk.Set
// @Failure 500 {object} response.ErrorResponse
// @Router /.well-known/openid-configuration/jwks [get]
func Jwks(cfg *configs.OAuth2, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 检查配置是否有效
		if cfg.Manager == nil {
			log.Error("OAuth2 manager configuration is missing")
			c.JSON(http.StatusInternalServerError, response.InternalServerError("server configuration error"))
			return
		}

		jwtPublicKeyPath := cfg.Manager.JWTPublicKey
		if jwtPublicKeyPath == "" {
			log.Error("JWT public key path not configured")
			c.JSON(http.StatusInternalServerError, response.InternalServerError("server configuration error"))
			return
		}

		pubKeyData, err := os.ReadFile(jwtPublicKeyPath)
		if err != nil {
			log.Error("Failed to read JWT public key", zap.String("path", jwtPublicKeyPath), zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.InternalServerError("could not read public key"))
			return
		}

		key, err := jwk.ParseKey(pubKeyData, jwk.WithPEM(true))
		if err != nil {
			log.Error("Failed to parse JWK", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.InternalServerError("could not parse public key"))
			return
		}

		// 设置或生成kid
		if kid := key.KeyID(); kid == "" {
			desiredKid := generateKeyID([]byte(cfg.Manager.Kid))
			if err := key.Set(jwk.KeyIDKey, desiredKid); err != nil {
				log.Error("Failed to set key ID", zap.Error(err))
			} else {
				log.Info("Generated kid for JWK", zap.String("kid", desiredKid))
			}
		}

		// 设置算法（重要）
		if alg := key.Algorithm(); alg == nil || alg.String() == "" {
			if err := key.Set(jwk.AlgorithmKey, "RS256"); err != nil {
				log.Warn("Failed to set algorithm for key", zap.Error(err))
			}
		}

		// 设置密钥用途
		if err := key.Set(jwk.KeyUsageKey, jwk.ForSignature); err != nil {
			log.Warn("Failed to set key usage", zap.Error(err))
		}

		set := jwk.NewSet()
		if err := set.AddKey(key); err != nil {
			log.Error("Failed to add key to JWKS", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.InternalServerError("could not build JWKS"))
			return
		}

		c.JSON(http.StatusOK, set)
	}
}

// generateKeyID 生成健壮且唯一的key ID
// 基于密钥内容生成 kid
func generateKeyID(key []byte) string {
	hash := sha256.Sum256(key)
	return base64.RawURLEncoding.EncodeToString(hash[:8])
}
