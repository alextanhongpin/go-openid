package code

type Service interface {
	Validate(code Code) error
	Code() Code
}
