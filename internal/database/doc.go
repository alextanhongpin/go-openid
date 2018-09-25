package database

// Contains the database implementation. The reason why the implementation is
// not placed together with the service (e.g. /service/user/repository.go) is
// because the storage might be accessed by another service. Hence, by
// separating them this allows another service to use the database
// independently. It is arguable that each service should have their own
// storage in microservice architecture/domain-driven-design etc, but there are
// exceptions. Hence, in a less strict design, this can still be allowed.
