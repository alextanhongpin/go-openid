package code

type Repository interface {
	WithCode(code string) (Code, error)
	Delete(code string) (bool, error)
	Create(code Code) (bool, error)
}
