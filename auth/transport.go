package auth

import (
  "context"
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "time"

  "github.com/alextanhongpin/go-openid/app"
  "github.com/dgrijalva/jwt-go"
  "github.com/julienschmidt/httprouter"
  "github.com/justinas/alice"
)

// FeatureToggle allows you to toggle the feature
func FeatureToggle(isEnabled bool) func(app.Env) {
  return func(env app.Env) {
    // Don't run this
    if !isEnabled {
      return
    }
    e := endpoint{userService{db: env.DB}}
    r := env.Router
    t := env.Tmpl

    // API
    r.GET("/api/users/:id", e.getUserHandler())
    r.GET("/api/users", e.getUsersHandler())

    // Static & Forms
    r.GET("/users/:id", e.viewUserHandler(t))

    r.POST("/login", e.loginHandler())
    r.GET("/login", e.loginViewHandler(t))

    r.GET("/register", e.registerViewHandler(t))
    r.POST("/register", e.registerHandler(injectJWT(encodeRegisterResponseJSON)))

    // Example middleware
    r.GET("/context", middleware(handler))

    r.GET("/middleware/:id",
      chain(alice.New(
        middlewareOne,
        middlewareTwo,
      ).ThenFunc(finalHandler)))

    r.GET("/cookie", chain(alice.New(mockCookie, isLoggedIn).ThenFunc(cookieHandler)))

    s := MakeAuthService(env.DB)
    se := MakeServerEndpoints(s)

    r.GET("/api/v2/users/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
      w.Write([]byte("Hello world" + ps.ByName("id")))
    })

    r.GET("/api/v3/users/:id", se.GetUser())
    // r.GET("/api/v3/users/:id", chain(alice.New(
    //   se.GetUser,
    // ).ThenFunc(finalHandler)))
  }
}

func middlewareOne(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // do stuff
    log.Println("At middleware one:start")
    next.ServeHTTP(w, r)
    log.Println("At middleware one:end")
  })
}

func middlewareTwo(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // do stuff
    log.Println("At middleware two:start")
    next.ServeHTTP(w, r)
    log.Println("At middleware two:end")
  })
}

// The last chain must be a http.HandlerFunc
func finalHandler(w http.ResponseWriter, r *http.Request) {
  params := r.Context().Value("params").(httprouter.Params)
  w.Write([]byte("Hello world" + params.ByName("id")))
}

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

type Claims struct {
  UserID string `json:"user_id"`
  jwt.StandardClaims
}

// Add jwt token to the response
func injectJWT(next httprouter.Handle) httprouter.Handle {
  return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

    userID := r.Context().Value("userID").(string)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
      userID,
      jwt.StandardClaims{
        ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
      },
    })

    tokenString, err := token.SignedString([]byte("$ecret"))
    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
    }
    expiration := time.Now().Add(365 * 24 * time.Hour)
    cookie := http.Cookie{Name: "access_token", Value: tokenString, Expires: expiration}
    http.SetCookie(w, &cookie)

    ctx := context.WithValue(r.Context(), "jwtToken", tokenString)
    next(w, r.WithContext(ctx), ps)
  })
}

type Key string

var userID Key

func validate(next httprouter.Handle) httprouter.Handle {
  return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    tokenString := "your_token_here"
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
      return []byte("$ecret"), nil
    })
    if err != nil {
      http.NotFound(w, r)
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
      ctx := context.WithValue(r.Context(), userID, *claims)
      next(w, r.WithContext(ctx), ps)
    } else {
      http.NotFound(w, r)
      return
    }
  })

}
func encodeRegisterResponseJSON(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  jwtToken := r.Context().Value("jwtToken").(string)
  userID := r.Context().Value("userID").(string)
  log.Printf("Found jwtToken:%s, and userID: %s", jwtToken, userID)
  response := registerResponse{
    OK:          true,
    RedirectURI: "/users/" + userID,
    AccessToken: jwtToken,
  }
  json.NewEncoder(w).Encode(response)

}
