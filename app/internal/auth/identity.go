package auth

import "context"

type Identity struct {
	UserID string   `json:"user_id"`
	Role   UserRole `json:"role"`
}

type ctxKey string

const identityKey ctxKey = "identity"

func NewIdentity(userID string, role UserRole) *Identity {
	return &Identity{
		UserID: userID,
		Role:   role,
	}
}

func FromContext(ctx context.Context) *Identity {
	identity := ctx.Value(identityKey)
	return identity.(*Identity)
}

func WithIdentity(ctx context.Context, id *Identity) context.Context {
	return context.WithValue(ctx, identityKey, id)
}

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleViewer UserRole = "viewer"
	UserRoleEditor UserRole = "editor"
)

func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleAdmin,
		UserRoleViewer,
		UserRoleEditor:
		return true
	}
	return false
}
