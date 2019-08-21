package validation

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

// Validation provides a path to a specific element in a struct and the problem that triggered the error
type Validation struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

// Error provides a message and a list of issues with the validation
type Error struct {
	Error  string       `json:"error"`
	Errors []Validation `json:"validation"`
}

// ValidateStruct validates a struct and return a slice of errors
func ValidateStruct(data interface{}) *Error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	err := validate.Struct(data)
	if err != nil {
		out := Error{
			Error: "validation error",
		}

		switch e := err.(type) {
		case *validator.InvalidValidationError:
			fmt.Println(e)
			out.Errors = append(out.Errors, Validation{
				Path:  "",
				Error: "invalid input",
			})
			return &out
		case validator.ValidationErrors:
			for _, err := range e {
				namespace := strings.Split(err.Namespace(), ".")
				namespace = namespace[1:]
				path := strings.Join(namespace, ".")

				out.Errors = append(out.Errors, Validation{
					Path:  path,
					Error: err.Tag(),
				})
			}
			return &out
		}
	}

	return nil
}
