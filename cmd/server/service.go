package main

type OpenIDService interface {
}

type Service struct {
	db *Database
}

func (s *Service) AuthorizationCodeFlow() {}

func (s *Service) TokenFlow() {}
func NewService(db *Database) *Service {
	if db == nil {
		db = NewDatabase()
	}
	return &Service{
		db: db,
	}
}
