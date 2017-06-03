package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/alextanhongpin/go-openid/app"

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
		err = e.svc.create(u)
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

func (e endpoint) loginViewHandler(tmpl *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tmpl.Render(w, "login", nil)
	}
}

func (e endpoint) viewUserHandler(tmpl *app.Template) httprouter.Handle {
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
		tmpl.Render(w, "login", user)
	}
}

// POST /register
// registerHandle register a user and return a redirect uri
func (e endpoint) registerHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var u User

		// Decode the body payload
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if the user exist
		user, err := e.svc.fetchOne(u.Email)

		// Error occured
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// No user found, create new account
		if user == nil {
			log.Println("registerHandler: No user found.")
			err := e.svc.create(u)
			if err != nil {
				log.Printf("registerHandler: Error creating user: %s", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			// Successfully created, return payload
			log.Println("registerHandler: successfully created user")
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")

			w.Header().Set("Access-Control-Expose-Headers", "Location")
			w.Header().Set("Location", "http://www.localhost:3000/prfile")
			w.Write([]byte(`{"success": true, "redirect_uri": "/users/` + u.Email + `"}`))
			return
		}
		log.Println("registerHandler: user exist")
		// w.WriteHeader(http.StatusUnauthorized)
		// w.Header().Set("Content-Type", "application/json")
		// w.Write([]byte(`{"message": "user already exists"}`))
		http.Error(w, `{"success": false, "message": "user already exists"}`, http.StatusUnauthorized)
	}
}

// GET register/
// registerViewHandler renders the register view
func (e endpoint) registerViewHandler(tmpl *app.Template) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tmpl.Render(w, "register", nil)
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
