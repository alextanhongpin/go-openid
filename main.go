package main

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

type Email struct {
	Gender string `validate:"eq=male|eq=female"` // checking enum
	Value  string `validate:"email,required"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func main() {

	validate = validator.New()
	a := Email{Gender: "malea", Value: "john.doe@mail.com"}
	err := validate.Struct(a)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}

}
