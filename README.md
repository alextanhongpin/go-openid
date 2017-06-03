# OpenID Connect with golang

Building an openid connect implementation with golang (WIP)

![Login Screen](./assets/login.png)


## Todos
- [ ] Reading Config
- [ ] Validating Model
- [ ] JWT
- [ ] Hashing Password


## Reading Config
## Validating Model
```go
import "gopkg.in/go-playground/validator.v9"

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func main() {
    validate = validator.New()
    a := Email{Gender: "malea", Value: "john.doe@mail.com"}
    err := validate.Struct(a)
    if err != nil {
        fmt.Println("error: " + err.Error())
    }
}
```

## Creating Middleware with julienschmidt httprouter.Handle

```go
func StdToJulienMiddleware(next http.Handler) httprouter.Handle {

    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        // do stuff
        next.ServeHTTP(w, r)
    }
}

// Pure "github.com/julienschmidt/httprouter" middleware
func JulienToJulienMiddleware(next httprouter.Handle) httprouter.Handle {

    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        // do stuff
        next(w, r, ps)
    }
}

func JulienHandler() httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        // do stuff
    }
}

func StdHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // do stuff
    })
}

func main() {
    router := httprouter.New()
    router.POST("/api/user/create", StdToJulienMiddleware(StdHandler()))
    router.GET("/api/user/create", JulienToJulienMiddleware(JulienHandler()))
    log.Fatal(http.ListenAndServe(":8000", router))
}
```

Simple middleware

```go
type key string

const ctxName key = "id"

func middleware(next httprouter.Handle) httprouter.Handle {
  return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    ctx := context.WithValue(r.Context(), ctxName, 12345)
    next(w, r.WithContext(ctx), ps)
  })
}
func handler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  reqID := r.Context().Value(ctxName).(int)
  fmt.Fprintf(w, "hello request id: %d", reqID)
}
```


### mgo session.Copy or session.Clone
First of all, we need to see the difference between mgo.Session.Copy() and mgo.Session.Clone(). While go.Session.Clone() returns a new session, the session uses the same socket connection. That isn't necessarily a bad thing, but keep in mind that on the server side, a stack is allocated per connection. So the sessions would share the same stack. Depending on your use cases, that may make a big difference.

And here is the problem – if you open a new socket connect for each record, this leads to a three way handshake, which is slowish. Reusing the same socket reduces this overhead, but there still is some and has the drawback described above.



### Cookie

Setting a cookie

```go
    import "net/http"

    expiration := time.Now().Add(365 * 24 * time.Hour)
    cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
    http.SetCookie(w, &cookie)
```


Getting a cookie
```go
    cookie, _ := r.Cookie("username")
    fmt.Fprint(w, cookie)
```


