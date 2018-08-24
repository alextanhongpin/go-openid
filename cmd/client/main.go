package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var ErrForbidden = errors.New("forbidden")

type Client struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	Get(id string) (*Client, error)
	Register(Client) error
	// Update()
	// Delete()
}

type clientService struct {
	sync.RWMutex
	db map[string]Client
}

func (c *clientService) Get(id string) (*Client, error) {
	c.RLock()
	defer c.RUnlock()
	if c, ok := c.db[id]; !ok {
		return nil, ErrForbidden
	} else {
		return &c, nil
	}
}

func (c *clientService) Register(req Client) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.db[req.ID]; ok {
		return ErrForbidden
	}
	c.db[req.ID] = req
	return nil
}

func main() {
	port := 8080

	s := &clientService{
		db: make(map[string]Client),
	}
	e := makeEndpoints(s)

	h := http.NewServeMux()
	h.HandleFunc("/connect/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			e.GetClient(w, r)
		case http.MethodPost:
			e.PostClientRegister(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
			fmt.Fprintf(w, "not implemented")
		}
	})

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("listening to port *:%d. press ctrl + c to cancel\n", port)
	log.Fatal(srv.ListenAndServe())
}

type Endpoint = http.HandlerFunc
type Endpoints struct {
	GetClient          Endpoint
	PostClientRegister Endpoint
}

func makeEndpoints(s Service) *Endpoints {
	return &Endpoints{
		GetClient:          makeGetClientEndpoint(s),
		PostClientRegister: makePostClientRegisterEndpoint(s),
	}
}

func makeGetClientEndpoint(s Service) Endpoint {
	return func(w http.ResponseWriter, r *http.Request) {
		cid := r.URL.Query().Get("client_id")
		if cid == "" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, ErrForbidden.Error())
			return
		}
		c, err := s.Get(cid)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, err.Error())
			return
		}
		// Check authorization header
		// Check if client exists in database
		// Return client
		// If not found, return 403, not 404
		json.NewEncoder(w).Encode(c)
	}
}

func makePostClientRegisterEndpoint(s Service) Endpoint {
	return func(w http.ResponseWriter, r *http.Request) {
		// if authHdr := r.Header.Get("Authorization"); authHdr == "" {
		// 	// Validate authorization header
		// }
		var req Client
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			return
		}

		if err := s.Register(req); err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, err.Error())
			return
		}
		// Validate request
		// Check if client already exist
		// Store to database
		// Return the response
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		fmt.Fprintf(w, "success")
	}
}
