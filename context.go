


package openid

import "context"

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

var (
	UserIDContextKey = ContextKey("user_id")
	AuthContextKey   = ContextKey("authorization")
)

func SetUserIDContextKey(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}

func GetUserIDContextKey(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	return userID, ok
}

func SetAuthContextKey(ctx context.Context, authorization string) context.Context {
	return context.WithValue(ctx, AuthContextKey, authorization)
}

func GetAuthContextKey(ctx context.Context) (string, bool) {
	auth, ok := ctx.Value(AuthContextKey).(string)
	return auth, ok
}
