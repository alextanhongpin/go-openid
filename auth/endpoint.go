package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	errorTemplateNotFound = errors.New("Template not found")
)

type endpoint struct {
	svc service
}

func (e endpoint) createUserHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		fmt.Println(u, err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(u)
		err = e.svc.create(u.Email, u.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `{"ok": true, "redirect_uri": "/users/%s"}`, u.Email)
	}
}

func (e endpoint) getUserHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := e.svc.fetchOne(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// No user found
		if user == nil {
			w.Write([]byte("no user found"))
			return
		}
		w.Write([]byte(user.Name))
	}
}

func (e endpoint) loginHandler(tmpl map[string]*template.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		t, ok := tmpl["login"]
		if !ok {
			http.Error(w, errorTemplateNotFound.Error(), http.StatusBadRequest)
			return
		}
		t.Execute(w, nil)
	}
}

func (e endpoint) viewUserHandler(tmpl map[string]*template.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		t, ok := tmpl["user"]
		if !ok {
			http.Error(w, errorTemplateNotFound.Error(), http.StatusBadRequest)
			return
		}
		user, err := e.svc.fetchOne(ps.ByName("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// No user found
		if user == nil {
			w.Write([]byte("no user found"))
			return
		}
		t.Execute(w, user)
	}
}

// func getUserHandler(db *memdb.MemDB) httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 		txn := db.Txn(false)
// 		defer txn.Abort()

// 		// Lookup by id
// 		// raw, err := txn.First("user", "id", "john.doe@mail.com")
// 		fmt.Println(r.FormValue("email"))

// 		result, err := txn.Get("user", "id", r.FormValue("email"))
// 		if err != nil {
// 			panic(err)
// 		}
// 		var users []string
// 		// for raw := result.Next(); raw != nil; {
// 		// 	users = append(users, raw.(*User).Name)
// 		// }
// 		for i := 0; i < 10; i++ {
// 			raw := result.Next()
// 			if raw == nil {
// 				break
// 			}
// 			users = append(users, raw.(*User).Name)
// 		}

// 		fmt.Println(len(users))
// 		w.Write([]byte(strings.Join(users, "-")))
// 	}
// }
