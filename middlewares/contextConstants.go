package middlewares

type contextKey int

// OriginContextKey will be used as the request context key where the client's IP is stored
const (
	OriginContextKey contextKey = iota
)
