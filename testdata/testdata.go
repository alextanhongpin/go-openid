package testdata

// Mocking should be done at higher-level for things that do not need to be
// tested, or infrastructure(database, message queue), as well as external
// services. The layer containing business logic should not be mocked, as that
// is the area we are primarily concerned with.
