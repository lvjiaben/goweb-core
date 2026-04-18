package rbac

import (
	"context"

	"github.com/lvjiaben/goweb-core/httpx"
)

const identityKey = "rbac.identity"

type Identity struct {
	UserID   int64
	UserType string
	RoleIDs  []int64
	IsSuper  bool
}

type PermissionChecker interface {
	HasPermission(ctx context.Context, identity *Identity, permissionCode string) (bool, error)
}

type MenuItem struct {
	ID             int64      `json:"id"`
	ParentID       int64      `json:"parent_id"`
	Name           string     `json:"name"`
	Title          string     `json:"title"`
	Path           string     `json:"path"`
	Component      string     `json:"component"`
	MenuType       string     `json:"menu_type"`
	Icon           string     `json:"icon"`
	Sort           int        `json:"sort"`
	PermissionCode string     `json:"permission_code"`
	Children       []MenuItem `json:"children,omitempty"`
}

type MenuProvider interface {
	GetMenus(ctx context.Context, identity *Identity) ([]MenuItem, error)
}

func SetIdentity(c *httpx.Context, identity *Identity) {
	c.Set(identityKey, identity)
}

func GetIdentity(c *httpx.Context) (*Identity, bool) {
	value, ok := c.Get(identityKey)
	if !ok {
		return nil, false
	}
	identity, ok := value.(*Identity)
	return identity, ok
}

func RequirePermission(checker PermissionChecker) httpx.Middleware {
	return func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(c *httpx.Context) {
			route := c.Route()
			if route == nil || route.PermissionCode == "" {
				next(c)
				return
			}

			identity, ok := GetIdentity(c)
			if !ok || identity == nil {
				c.Unauthorized("unauthorized")
				return
			}

			allowed, err := checker.HasPermission(c.Request.Context(), identity, route.PermissionCode)
			if err != nil {
				c.Error(err)
				return
			}
			if !allowed {
				c.Forbidden("permission denied")
				return
			}
			next(c)
		}
	}
}
