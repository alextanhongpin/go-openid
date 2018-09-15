package clientsvc

import (
	"log"
	"time"

	"github.com/alextanhongpin/go-openid/app"
	"github.com/alextanhongpin/go-openid/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GetClientMetadata 	Get a client metadata by id
// GetClientsMetadata   Get a list of clients metadata with pagination

type Service interface {
	GetClientMetadata(getClientMetadataRequest) (*getClientMetadataResponse, error)
	GetClientsMetadata(getClientsMetadataRequest) (*getClientsMetadataResponse, error)
	PostClientMetadata(postClientMetadataRequest) (*postClientMetadataResponse, error)
	UpdateClientMetadata(updateClientMetadataRequest) (*updateClientMetadataResponse, error)
	DeleteClientMetadata(deleteClientMetadataRequest) (*deleteClientMetadataResponse, error)
}

type clientsvc struct {
	db *app.Database
}

// MakeClientService creates a new client service that can access the client repository
func MakeClientService(db *app.Database) Service {
	return &clientsvc{db}
}

func (s clientsvc) GetClientMetadata(req getClientMetadataRequest) (*getClientMetadataResponse, error) {
	session := s.db.NewSession()
	defer session.Close()

	var client models.Client

	c := s.db.Collection("client", session)
	err := c.FindId(bson.ObjectIdHex(req.ID)).One(&client)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	log.Printf("GetClientMetadata type=service client_metadata=%v", client.ClientMetadata)
	log.Printf("GetClientMetadata type=service client_id=%s client_secret=%s", client.ClientID, client.ClientSecret)
	return &getClientMetadataResponse{
		Data: client,
	}, nil
}

func (s clientsvc) GetClientsMetadata(req getClientsMetadataRequest) (*getClientsMetadataResponse, error) {
	session := s.db.NewSession()
	defer session.Close()
	c := s.db.Collection("client", session)

	var clients []models.Client
	err := c.Find(bson.M{}).All(&clients)
	if err != nil {
		return nil, err
	}

	// TODO: Remove client id and secrets from each clients

	// iter := c.Find(bson.M{}).Iter()
	// client := ClientMetadata{}
	// for iter.Next(&client) {
	// 	clientMetadatas = append(clientMetadatas, client)
	// }

	return &getClientsMetadataResponse{
		Data: clients,
	}, nil
}

func (s clientsvc) PostClientMetadata(req postClientMetadataRequest) (*postClientMetadataResponse, error) {
	session := s.db.NewSession()
	defer session.Close()

	c := s.db.Collection("client", session)

	timestamp := time.Now().Unix()
	expireTimestamp := time.Now().Add(1 * time.Hour).Unix()
	// req.ClientMetadata.CreatedAt = &presentTime
	// req.ClientMetadata.UpdatedAt = &presentTime

	// Generate client secret and client id

	// go v1.8 enables us to carry out conversion rules from one struct to another to keep code DRY
	// client := ClientMetadata(req.ClientMetadata)
	client := models.Client{
		ID:                    bson.NewObjectId(),
		IsPublished:           false,
		Version:               "0.0.1",
		CreatedAt:             &timestamp,
		UpdatedAt:             &timestamp,
		ClientID:              "generated_client_id",
		ClientSecret:          "generated_client_secret",
		ClientIDIssuedAt:      &timestamp,
		ClientSecretExpiresAt: &expireTimestamp,
		ClientMetadata:        req.ClientMetadata,
		// RegistrationAccessToken : "",
		// RegistrationClientURI  : "",
	}
	err := c.Insert(&client)
	if err != nil {
		return nil, err
	}
	// TODO: Return the correct model
	return &postClientMetadataResponse{
		Data: req.ClientMetadata,
	}, nil
}

func (s clientsvc) UpdateClientMetadata(req updateClientMetadataRequest) (*updateClientMetadataResponse, error) {
	session := s.db.NewSession()
	defer session.Close()

	c := s.db.Collection("client", session)

	log.Printf("UpdateClientMetadata type=service req=%v \n", req)
	// id, err := utils.ValidateId(req.ID.Hex())
	// if err != nil {
	// 	return &res, err
	// }
	err := c.Update(bson.M{"_id": req.ID}, bson.M{"$set": bson.M{
		"client_metadata.client_name":      req.ClientName,
		"client_metadata.application_type": req.ApplicationType,
		"client_metadata.redirect_uris":    req.RedirectURIs,
		"updated_at":                       time.Now(),
	}})
	if err != nil {
		log.Printf("UpdateClientMetadata type=service error=%v \n", err)
		return nil, err
	}
	return &updateClientMetadataResponse{Ok: true}, nil
}

func (s clientsvc) DeleteClientMetadata(req deleteClientMetadataRequest) (*deleteClientMetadataResponse, error) {
	session := s.db.NewSession()
	defer session.Close()

	c := s.db.Collection("client", session)

	log.Printf("DeleteClientMetadata type=service event=delete_client params=%v\n", req)
	err := c.RemoveId(req.ID)
	if err != nil {
		log.Printf("DeleteClientMetadata type=service event=error_deleting_client error=%v\n", err)
		return nil, err
	}

	log.Println("DeleteClientMetadata type=service event=success_delete_client")
	return &deleteClientMetadataResponse{Ok: true}, nil
}
