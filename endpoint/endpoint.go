package endpoint

// Endpoint is the fundamental building block of servers and clients.
type Endpoint func(request interface{}) (response interface{}, err error)

// Middleware is a chainable behavior modifier for endpoints.
type Middleware func(Endpoint) Endpoint

// Chain is a helper function for composing middlewares.
func Chain(outer Middleware, others ...Middleware) Middleware {
  return func(next Endpoint) Endpoint {
    for i := len(others) - 1; i >= 0; i-- {
      next = others[i](next)
    }
    return outer(next)
  }
}
