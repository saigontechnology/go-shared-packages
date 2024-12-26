package authmiddleware

import (
	"sync"
)

type ClaimData struct {
	Claim Claim `json:"data"`
}

type Claim struct {
	AccountID    string   `json:"account_id"`
	AccountRoles []string `json:"account_roles"`
}

var (
	providerOnce     sync.Once
	providerInstance *provider
)

type AuthMiddleware interface {
	CheckPermission(token, path, method string) (*Claim, error)
}

type Provider interface {
	GetMiddleware() AuthMiddleware
}

type provider struct{}

// GetProvider singleton implementation makes sure only one Provider is created to avoid duplicated database connection pools.
func GetProvider() Provider {
	providerOnce.Do(func() {
		providerInstance = &provider{}
	})

	return providerInstance
}

func (p *provider) GetMiddleware() AuthMiddleware {
	return NewAuthMiddleware()
}
