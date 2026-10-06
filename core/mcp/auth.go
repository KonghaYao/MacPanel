package mcp

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"time"

	appauth "github.com/1Panel-dev/1Panel/core/app/auth"
	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/app/repo"
	"github.com/1Panel-dev/1Panel/core/constant"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/1Panel-dev/1Panel/core/init/session/psession"
	"github.com/1Panel-dev/1Panel/core/utils/encrypt"
	"github.com/gin-gonic/gin"
	mcpauth "github.com/modelcontextprotocol/go-sdk/auth"
	"gorm.io/gorm"
)

const authenticatedAPIKeyExtraKey = "macpanel_api_key"

// AuthenticatedAPIKey holds the resolved Panel API key for an MCP request.
type AuthenticatedAPIKey struct {
	KeyID    string
	KeyName  string
	Kind     string // "apiKey" or "legacy"
	Secret   string
	Owner    appauth.APIKeyOwner
	ClientIP string
	Revision uint64
}

// APIKeyFromContext returns the authenticated API key stored by MCP bearer auth.
func APIKeyFromContext(ctx context.Context) (*AuthenticatedAPIKey, bool) {
	if info := mcpauth.TokenInfoFromContext(ctx); info != nil {
		return APIKeyFromTokenInfo(info)
	}
	return nil, false
}

// APIKeyFromTokenInfo extracts the authenticated API key from MCP TokenInfo.
func APIKeyFromTokenInfo(info *mcpauth.TokenInfo) (*AuthenticatedAPIKey, bool) {
	if info == nil || info.Extra == nil {
		return nil, false
	}
	value, ok := info.Extra[authenticatedAPIKeyExtraKey]
	if !ok {
		return nil, false
	}
	key, ok := value.(AuthenticatedAPIKey)
	if !ok {
		return nil, false
	}
	return &key, true
}

// VerifyBearerAPIKey validates MCP Authorization bearer tokens against Panel API keys.
func VerifyBearerAPIKey(ctx context.Context, token string, req *http.Request) (*mcpauth.TokenInfo, error) {
	if token == "" {
		return nil, fmt.Errorf("%w: missing bearer token", mcpauth.ErrInvalidToken)
	}

	enabled, err := apiInterfaceEnabled()
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, invalidBearerToken(errors.New("API interface disabled"))
	}

	authKey, err := resolveBearerAPIKey(token, req)
	if err != nil {
		return nil, invalidBearerToken(err)
	}

	tokenInfo := &mcpauth.TokenInfo{
		UserID: authKey.Owner.ID,
		Extra: map[string]any{
			authenticatedAPIKeyExtraKey: *authKey,
		},
	}
	if expiresAt := authKeyExpiresAt(authKey); expiresAt != nil {
		tokenInfo.Expiration = *expiresAt
	}
	return tokenInfo, nil
}

func apiInterfaceEnabled() (bool, error) {
	status, err := repo.NewISettingRepo().GetValueByKey("ApiInterfaceStatus")
	if err != nil {
		return false, err
	}
	return status == constant.StatusEnable, nil
}

func resolveBearerAPIKey(token string, req *http.Request) (*AuthenticatedAPIKey, error) {
	if authKey, err := matchManagedAPIKey(token, req); err != nil {
		return nil, err
	} else if authKey != nil {
		return authKey, nil
	}
	return matchLegacyAPIKey(token, req)
}

func matchManagedAPIKey(token string, req *http.Request) (*AuthenticatedAPIKey, error) {
	master, err := encrypt.APIKeyEncryptionKey()
	if err != nil {
		return nil, err
	}

	provider := appauth.PanelAPIKeyOwnerProvider{}
	now := time.Now()
	var candidates []model.APIKey
	var matched *model.APIKey

	err = global.DB.Model(&model.APIKey{}).
		Where("status <> ?", "Revoked").
		Order("id").
		FindInBatches(&candidates, 128, func(_ *gorm.DB, _ int) error {
			for _, key := range candidates {
				if key.SecretVersion != 1 {
					continue
				}
				secret, err := encrypt.DecryptAPIKey(key.SecretCiphertext, key.ID, master)
				if err != nil {
					continue
				}
				if subtle.ConstantTimeCompare([]byte(secret), []byte(token)) != 1 {
					continue
				}
				copy := key
				matched = &copy
				return errBearerAPIKeyMatched
			}
			return nil
		}).Error
	if err != nil && !errors.Is(err, errBearerAPIKeyMatched) {
		return nil, err
	}
	if matched == nil {
		return nil, nil
	}

	key, err := repo.NewAPIKeyRepo().Get(matched.ID)
	if err != nil || key.Revision != matched.Revision {
		return nil, errors.New("API key status changed")
	}
	if err = appauth.ValidateAPIKeyState(key, now); err != nil {
		return nil, err
	}
	owner, err := provider.ResolveAPIKeyOwner(nil, key.OwnerType, key.OwnerID)
	if err != nil {
		return nil, err
	}
	secret, err := appauth.APIKeySecret(key)
	if err != nil {
		return nil, err
	}
	clientIP := clientIPFromRequest(req, key.APITrustedProxies)
	if !appauth.IsIPInWhiteList(clientIP, key.IPWhiteList) {
		return nil, errors.New("client IP not allowed")
	}

	return &AuthenticatedAPIKey{
		KeyID:    key.ID,
		KeyName:  key.Name,
		Kind:     "apiKey",
		Secret:   secret,
		Owner:    owner,
		ClientIP: clientIP,
		Revision: key.Revision,
	}, nil
}

func matchLegacyAPIKey(token string, req *http.Request) (*AuthenticatedAPIKey, error) {
	config, err := appauth.LoadAPIAuthConfig(nil)
	if err != nil {
		return nil, err
	}
	if config.ApiInterfaceStatus != constant.StatusEnable || config.ApiKey == "" {
		return nil, errors.New("legacy API key disabled")
	}
	if subtle.ConstantTimeCompare([]byte(config.ApiKey), []byte(token)) != 1 {
		return nil, errors.New("API key not found")
	}

	provider := appauth.PanelAPIKeyOwnerProvider{}
	owner, err := provider.ResolveAPIKeyOwner(nil, appauth.APIKeyOwnerPanelAdmin, psession.SuperAdminSessionUserID)
	if err != nil {
		return nil, err
	}
	clientIP := clientIPFromRequest(req, config.ApiTrustedProxies)
	if !appauth.IsIPInWhiteList(clientIP, config.IpWhiteList) {
		return nil, errors.New("client IP not allowed")
	}

	return &AuthenticatedAPIKey{
		KeyID:    "legacy",
		KeyName:  "Legacy API Key",
		Kind:     "legacy",
		Secret:   config.ApiKey,
		Owner:    owner,
		ClientIP: clientIP,
	}, nil
}

func clientIPFromRequest(req *http.Request, trustedProxies string) string {
	if req == nil {
		return ""
	}
	c := &gin.Context{Request: req}
	return appauth.GetAPIClientIP(c, trustedProxies)
}

func authKeyExpiresAt(key *AuthenticatedAPIKey) *time.Time {
	if key == nil || key.Kind != "apiKey" || key.KeyID == "" {
		return nil
	}
	record, err := repo.NewAPIKeyRepo().Get(key.KeyID)
	if err != nil || record.ExpiresAt == nil {
		return nil
	}
	expiresAt := *record.ExpiresAt
	return &expiresAt
}

var errBearerAPIKeyMatched = errors.New("bearer API key matched")

func invalidBearerToken(cause error) error {
	if cause != nil && global.LOG != nil {
		global.LOG.Debugf("mcp bearer auth failed: %v", cause)
	}
	return mcpauth.ErrInvalidToken
}
