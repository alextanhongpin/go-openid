package client

import "context"

type contextKey string

type Context struct {
}

const (
	clientSecretKey contextKey = "client_secret"
	clientIDKey     contextKey = "client_id"
)

func (c Context) WithClientSecret(ctx context.Context, clientSecret string) context.Context {
	return context.WithValue(ctx, clientSecretKey, clientSecret)

}
func ClientSecret(ctx context.Context) string {
	value, ok := ctx.Value(clientSecretKey)
	if !ok {
		return ""
	}
	return value.(string)
}
