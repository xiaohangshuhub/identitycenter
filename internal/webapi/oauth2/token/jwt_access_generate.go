package token

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/xiaohangshuhub/xiaohangshu/configs"
)

// CustomJWTAccessClaims 自定义JWT AccessClaims
type CustomJWTAccessClaims struct {
	jwt.RegisteredClaims
}

// Valid claims verification
func (a *CustomJWTAccessClaims) Valid() error {
	if a.ExpiresAt != nil && time.Unix(a.ExpiresAt.Unix(), 0).Before(time.Now()) {
		return errors.New("invalid access token")
	}
	return nil
}

// CustomJWTAccessGenerate 自定义JWT AccessGenerate
type CustomJWTAccessGenerate struct {
	SignedKeyID  string
	SignedKey    []byte
	SignedMethod jwt.SigningMethod
	Issuer       string
}

// Token 生成访问令牌和刷新令牌
func (a *CustomJWTAccessGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	claims := &CustomJWTAccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{data.Client.GetID()},
			Subject:   data.UserID,
			ExpiresAt: jwt.NewNumericDate(data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn())),
			Issuer:    a.Issuer,
		},
	}

	token := jwt.NewWithClaims(a.SignedMethod, claims)
	if a.SignedKeyID != "" {
		token.Header["kid"] = a.SignedKeyID
	}
	var key interface{}
	if a.isEs() {
		v, err := jwt.ParseECPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", "", err
		}
		key = v
	} else if a.isRsOrPS() {
		v, err := jwt.ParseRSAPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", "", err
		}
		key = v
	} else if a.isHs() {
		key = a.SignedKey
	} else if a.isEd() {
		v, err := jwt.ParseEdPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", "", err
		}
		key = v
	} else {
		return "", "", errors.New("unsupported sign method")
	}

	access, err := token.SignedString(key)
	if err != nil {
		return "", "", err
	}
	refresh := ""

	if isGenRefresh {
		t := uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).String()
		refresh = base64.URLEncoding.EncodeToString([]byte(t))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return access, refresh, nil
}

// 签名算法判断
func (a *CustomJWTAccessGenerate) isEs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "ES")
}

func (a *CustomJWTAccessGenerate) isRsOrPS() bool {
	isRs := strings.HasPrefix(a.SignedMethod.Alg(), "RS")
	isPs := strings.HasPrefix(a.SignedMethod.Alg(), "PS")
	return isRs || isPs
}

func (a *CustomJWTAccessGenerate) isHs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "HS")
}

func (a *CustomJWTAccessGenerate) isEd() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "Ed")
}

// NewConsumtJWTAccessGenerate 创建JWT AccessGenerate
func NewCustomJWTAccessGenerate(cfg *configs.OAuth2) oauth2.AccessGenerate {
	keyBytes, err := os.ReadFile(cfg.Manager.JWTPrivateKey)
	if err != nil {
		panic(err)
	}

	// 根据配置选择签名算法
	var signingMethod jwt.SigningMethod

	switch cfg.Manager.SigningMethod {
	case "RS256":
		signingMethod = jwt.SigningMethodRS256
	case "ES256":
		signingMethod = jwt.SigningMethodES256
	case "RS512":
		signingMethod = jwt.SigningMethodRS512
	default:
		signingMethod = jwt.SigningMethodRS256
	}
	// 生成key ID
	kid := generateKeyID([]byte(cfg.Manager.Kid))

	return &CustomJWTAccessGenerate{
		SignedKeyID:  kid,
		SignedKey:    keyBytes,
		SignedMethod: signingMethod,
		Issuer:       cfg.Issuer,
	}
}

// generateKeyID 生成健壮且唯一的key ID
// 基于密钥内容生成 kid
func generateKeyID(key []byte) string {
	hash := sha256.Sum256(key)
	return base64.RawURLEncoding.EncodeToString(hash[:8])
}
