package authmiddleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saigontechnology/go-shared-packages/logger"
	"github.com/saigontechnology/go-shared-packages/must"
)

const (
	tokenKeyHeader = "Authorization"
	AccountIDKey   = "account_id"
)

var (
	ErrUnauthorized         = errors.New("unauthorized")
	ErrNoTokenProvided      = errors.New("no token provided")
	ErrInvalidPathAndMethod = errors.New("invalid path and method")
)

type authMiddleware struct {
	identityServer string
}

func (a authMiddleware) CheckPermission(token, path, method string) (*Claim, error) {

	client := &http.Client{Timeout: 5 * time.Second}

	authRequest := map[string]interface{}{
		"path":   path,
		"method": method,
	}
	requestBody, err := json.Marshal(authRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/identity/auth/v1/check", a.identityServer),
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrUnauthorized
	}
	data := ClaimData{}
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return nil, ErrUnauthorized
	}

	return &data.Claim, nil
}

func NewAuthMiddleware() AuthMiddleware {
	cfg, err := newConfig()
	must.NotFail(err)

	return &authMiddleware{
		identityServer: cfg.IdentityServer,
	}
}

type CheckPermissionRequestPayload struct {
	Path   string `json:"path" binding:"required"`
	Method string `json:"method" binding:"required"`
}

func AuthInjectionMiddleware() gin.HandlerFunc {

	authMiddleware := NewAuthMiddleware()

	return func(c *gin.Context) {
		var token string
		if tokenAuthHeader := c.GetHeader(tokenKeyHeader); tokenAuthHeader != "" {
			token = tokenAuthHeader
		}
		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		request := &CheckPermissionRequestPayload{
			Path:   c.Request.URL.Path,
			Method: c.Request.Method,
		}

		claim, err := authMiddleware.CheckPermission(token, request.Path, request.Method)

		if err != nil {
			logger.GetProvider().Logger().Error(c, err.Error())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(AccountIDKey, claim.AccountID)

		c.Next()
	}
}
