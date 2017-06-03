package app

// import (
//   "context"
//   "github.com/uber-go/zap"
// )

// type correlationIdType int

// const (
//   requestIDKey correlationIDType = iota
//   sessionIDKey
// )

// var logger zap.Logger

// func init() {
//   logger = zap.New(
//     zap.NewJSONEncoder(zap.TimeFormatter(TimestampField)),
//     zap.Fields(zap.Int("pid", os.Getpid()),
//       zap.String("exe", path.Base(os.Args[0]))),
//   )
// }

// // WithReqID returns a context which knows its request ID
// func WithReqID(ctx context.Context, reqID string) context.Context {
//   return context.WithValue(ctx, requestIDKey, reqID)
// }

// // WithSessionID returns a context which knows its session ID
// func WithSessionID(ctx context.Context, sessionID string) context.Context {
//   return context.WithValue(ctx, sessionIDKey, sessionID)
// }

// // Logger returns a zap logger with as much context as possible
// func Logger(ctx context.Context) zap.Logger {
//   newLogger := logger
//   if ctx != nil {
//     if ctxReqID, ok := ctx.Value(requestIDKey).(string); ok {
//       newLogger = newLogger.With(zap.String("reqID", ctxReqID))
//     }
//     if ctxSessionID, ok := ctx.Value(sessionIDKey).(string); ok {
//       newLogger = newLogger.With(zap.String("sessionID", ctxSessionID))
//     }
//   }
//   return newLogger
// }
