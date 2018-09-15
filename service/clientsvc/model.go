package clientsvc

import (
	"github.com/alextanhongpin/go-openid/models"
	"gopkg.in/mgo.v2/bson"
)

type getClientMetadataRequest struct {
	ID string `json:"id"`
}

type getClientMetadataResponse struct {
	Data models.Client `json:"data,omitempty"`
}

type getClientsMetadataRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
type getClientsMetadataResponse struct {
	Data  []models.Client `json:"data,omitempty"`
	Count int             `json:"count"`
}

type postClientMetadataRequest struct {
	models.ClientMetadata
}
type postClientMetadataResponse struct {
	// ID string `json:"id"`
	Data models.ClientMetadata `json:"data,omitempty"`
}

type updateClientMetadataRequest struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	models.ClientMetadata
}

type updateClientMetadataResponse struct {
	Ok bool `json:"ok"`
}

type deleteClientMetadataRequest struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
}

type deleteClientMetadataResponse struct {
	Ok bool `json:"ok"`
}
