package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alextanhongpin/go-openid/entity"
	"github.com/alextanhongpin/go-openid/pkg/session"
	"github.com/julienschmidt/httprouter"
)

// SessionAuthorized checks if the session is valid or not, and takes a boolean
// delegate. If true, it will pass the state to the next handler through
// context, if false it will directly return an error.
func SessionAuthorized(sessMgr session.Manager, h httprouter.Handle, delegate bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		sess, err := sessMgr.GetSession(r)
		if err != nil && !delegate {
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(map[string]interface {
			}{

				"error": err.Error(),
			}); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			return
		}
		// Pass the session down to the next handler only if it exists.
		if sess != nil {
			ctx := r.Context()
			ctx = context.WithValue(ctx, entity.ContextKeySession, sess)
			r = r.WithContext(ctx)
		}
		h(w, r, ps)
	}
}

// RedirectIfSessionExist redirects the user to targetURL if the session
// already exists.
// Usage: RedirectIfSessionExist(GetLogin, sessMgr, "/home")
func RedirectIfSessionExists(h httprouter.Handle, sessMgr *session.Manager, targetURL string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if exist := sessMgr.HasSession(r); exist {
			http.Redirect(w, r, targetURL, http.StatusFound)
			return
		}
		h(w, r, ps)
	}
}
