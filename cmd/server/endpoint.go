package main

// Endpoints represent the endpoints for the OpenIDConnect.
type Endpoints struct {
	service Service
}

// NewEndpoints returns a pointer to new endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		service: s,
	}
}

//
// // Authorize performs the authorization logic.
// func (e *Endpoints) Authorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//         // Construct request parameters
//         var req oidc.AuthorizationRequest
//         if err := qs.Decode(&req, r.URL.Query()); err != nil {
//                 http.Error(w, err.Error(), http.StatusForbidden)
//                 return
//         }
//
//         if !govalidator.IsURL(req.RedirectURI) {
//                 http.Error(w, oidc.ErrInvalidRedirectURI.Error(), http.StatusForbidden)
//                 return
//         }
//
//         // Prepare redirect uri
//         u, err := url.Parse(req.RedirectURI)
//         if err != nil {
//                 http.Error(w, err.Error(), http.StatusForbidden)
//                 return
//         }
//
//         // Call service
//         res, authErr := e.service.Authorize(r.Context(), &req)
//         if authErr != nil {
//                 u.RawQuery = qs.Encode(authErr).Encode()
//                 http.Redirect(w, r, u.String(), http.StatusFound)
//                 return
//         }
//
//         u.RawQuery = qs.Encode(res).Encode()
//         http.Redirect(w, r, u.String(), http.StatusFound)
// }
//
// // Token represents the token service.
// func (e *Endpoints) Token(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//         r.ParseForm()
//
//         w.Header().Set("Content-Type", "application/json")
//         w.Header().Set("Cache-Control", "no-store")
//         w.Header().Set("Pragma", "no-cache")
//
//         token, err := authheader.Basic(r.Header.Get("Authorization"))
//         if err != nil {
//                 w.WriteHeader(http.StatusForbidden)
//                 json.NewEncoder(w).Encode(oidc.ErrInvalidRequest)
//                 return
//         }
//
//         clientID, clientSecret := authheader.DecodeBase64(token)
//
//         var req oidc.AccessTokenRequest
//         if err := qs.Decode(&req, r.Form); err != nil {
//                 w.WriteHeader(http.StatusForbidden)
//                 json.NewEncoder(w).Encode(oidc.ErrInvalidRequest)
//                 return
//         }
//
//         res, err := e.service.Token(r.Context(), &req)
//         if err != nil {
//                 w.WriteHeader(http.StatusForbidden)
//                 json.NewEncoder(w).Encode(err)
//                 return
//         }
//
//         w.WriteHeader(http.StatusOK)
//         json.NewEncoder(w).Encode(res)
// }
//
// // RegisterClient represents the endpoint for client registration.
// func (e *Endpoints) RegisterClient(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//
//         w.Header().Set("Content-Type", "application/json")
//         w.Header().Set("Cache-Control", "no-store")
//         w.Header().Set("Pragma", "no-cache")
//
//         token, err := authheader.Bearer(r.Header.Get("Authorization"))
//         if err != nil {
//                 w.WriteHeader(http.StatusBadRequest)
//                 json.NewEncoder(w).Encode(oidc.ErrInvalidRequest)
//                 return
//         }
//
//         // TODO: Validate token.
//
//         // Check for authorization headers to see if the client can register
//         var req oidc.Client
//         if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//                 w.WriteHeader(http.StatusBadRequest)
//                 json.NewEncoder(w).Encode(oidc.ErrInvalidRequest)
//                 return
//         }
//
//         res, err := e.service.RegisterClient(r.Context(), &req)
//         if err != nil {
//                 w.WriteHeader(http.StatusBadRequest)
//                 json.NewEncoder(w).Encode(oidc.ErrInvalidRequest)
//                 return
//         }
//
// Set the appropriate headers
//         w.WriteHeader(http.StatusCreated)
//         json.NewEncoder(w).Encode(res)
// }
//
// // Client returns the authorized client information.
// func (e *Endpoints) Client(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//
//         w.Header().Set("Content-Type", "application/json")
//         w.Header().Set("Cache-Control", "no-store")
//         w.Header().Set("Pragma", "no-cache")
//
//         token, err := authheader.Bearer(r.Header.Get("Authorization"))
//         if err != nil {
//                 w.WriteHeader(http.StatusUnauthorized)
//                 json.NewEncoder(w).Encode(oidc.ErrInvalidRequest)
//                 return
//         }
//
//         // TODO: Validate token.
//
//         id := r.URL.Query().Get("client_id")
//         client, err := e.service.ReadClient(r.Context(), id)
//         if err != nil {
//                 w.WriteHeader(http.StatusUnauthorized)
//                 json.NewEncoder(w).Encode(err)
//                 return
//         }
//         json.NewEncoder(w).Encode(client)
// }
//
// // .well-known/webfinger
// func (e *Endpoints) Webfinger() {}
//
// // .well-known/openid-configuration
// func (e *Endpoints) Configuration() {}
//
// func (e *Endpoints) Authenticate(ctx context.Context, req *oidc.AuthenticationRequest) (*oidc.AuthenticationResponse, error) {
//         return nil, nil
// }
//
// // RefreshToken returns a new refresh token alongside with the id token.
// func (e *Endpoints) RefreshToken() {}
//
// func (e *Endpoints) UserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//
//         _, err := authheader.Bearer(r.Header.Get("Authorization"))
//         if err != nil {
//                 // Error
//                 err := oidc.ErrUnauthorizedClient
//                 msg := fmt.Sprintf(`error="%s" error_description="%s"`, err.Error(), "The access token expired")
//                 w.Header().Set("WWW-Authenticate", msg)
//                 http.Error(w, err.Error(), http.StatusForbidden)
//                 return
//         }
//
//         res, err := e.service.UserInfo(r.Context(), "")
//         if err != nil {
//                 w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
//                 http.Error(w, err.Error(), http.StatusForbidden)
//                 return
//         }
//
//         json.NewEncoder(w).Encode(res)
// }
