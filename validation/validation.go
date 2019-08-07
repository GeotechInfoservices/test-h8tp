package validation

import (
	"context"
	"encoding/json"
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

		for _, err := range err.(validator.ValidationErrors) {
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

	return nil
}

// Validate middleware.
// Unmarshals and validates a struct
func Validate(handlerFunc interface{}) Handler {
	handler := reflect.ValueOf(handlerFunc)
	handlerType := reflect.TypeOf(handlerFunc)

	return func(ctx context.Context, req request.Request) (response.Response, error) {
		// construct arguments
		var args []reflect.Value
		if handlerType.NumIn() != 3 {
			//return some error handling
		}
			eventType := handlerType.In(handlerType.NumIn() - 1)
			event := reflect.New(eventType)

			if err := json.Unmarshal(payload, event.Interface()); err != nil {
				return nil, err
			}
			args = append(args, event.Elem())

		response := handler.Call(args)

		// convert return values into (interface{}, error)
		var err error
		if len(response) > 0 {
			if errVal, ok := response[len(response)-1].Interface().(error); ok {
				err = errVal
			}
		}
		var val interface{}
		if len(response) > 1 {
			val = response[0].Interface()

		}

		return val, err
	}
}
