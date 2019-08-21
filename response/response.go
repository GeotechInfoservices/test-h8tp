package response

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
)

// Response use APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// InternalServerError response for API.
func InternalServerError() (Response, error) {

	resp := map[string]interface{}{
		"message": "internal server error",
	}

	// ignore error, since we know this will always pass given a simple string.
	b, _ := json.Marshal(resp)
	return Response{
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Origin":      "*",
		},
		StatusCode: 500,
		Body:       string(b),
	}, nil
}

// NotFound response for API.
func NotFound(msg string) (Response, error) {

	resp := map[string]interface{}{
		"message": msg,
	}

	// ignore error, since we know this will always pass given a simple string.
	b, _ := json.Marshal(resp)
	return Response{
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Origin":      "*",
		},
		StatusCode: 404,
		Body:       string(b),
	}, nil
}

// InvalidRequest response for API. This automatically wraps the error message in the correct format.
func InvalidRequest(msg string) (Response, error) {

	resp := map[string]interface{}{
		"message": msg,
	}

	// ignore error, since we know this will always pass given a simple string.
	b, _ := json.Marshal(resp)
	return Response{
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Origin":      "*",
		},
		StatusCode: 400,
		Body:       string(b),
	}, nil
}

// Unauthorized response for API. This automatically wraps the error message in the correct format.
func Unauthorized(msg string) (Response, error) {

	resp := map[string]interface{}{
		"message": msg,
	}

	// ignore error, since we know this will always pass given a simple string.
	b, _ := json.Marshal(resp)
	return Response{
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Origin":      "*",
		},
		StatusCode: 401,
		Body:       string(b),
	}, nil
}

// OK 200 response from api. This automatically marshals a struct and converts it to json
func OK(body interface{}) (Response, error) {
	resp := map[string]interface{}{
		"data": body,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		msg := map[string]string{"message": "error while retrieving items"}
		b, _ := json.Marshal(msg)
		return Response{
			StatusCode: 500,
			Body:       string(b),
		}, nil
	}

	return Response{
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Origin":      "*",
		},
		StatusCode: 200,
		Body:       string(b),
	}, nil
}

// Error for validation errors.
type Error struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

// ValidationError contains path for each error to be associated with
type ValidationError struct {
	Message string  `json:"error"`
	Errors  []Error `json:"validation"`
}

func (e ValidationError) Error() string {
	return e.Message
}

// HandleValidationError recursively walks a tree of errors and flattens them into json path style
func HandleValidationError(path string, errors error) []Error {
	var errs []Error

	switch err := errors.(type) {
	case govalidator.Errors:
		for _, e := range err.Errors() {
			errs = append(errs, HandleValidationError(path, e)...)
		}
	case govalidator.Error:
		path += "."
		path += err.Name

		return []Error{
			{
				Path:  path,
				Error: err.Error(),
			},
		}
	default:
		return errs
	}

	return errs
}

// BadInput implies an error in the input, according to the entity validation rules.
func BadInput(errors error) (Response, error) {
	out := ValidationError{
		Message: "invalid input",
	}
	switch err := errors.(type) {
	case govalidator.Errors:
		out.Errors = HandleValidationError("$", err)
	default:
		logrus.Errorf("Error while handling Bad Input. %+v", err)
		e, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})

		return Response{
			Headers: map[string]string{
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Allow-Origin":      "*",
			},
			StatusCode: 500,
			Body:       string(e),
		}, nil

	}

	raw, err := json.Marshal(out)
	if err != nil {
		// TODO: Do some checking
	}

	return Response{
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Origin":      "*",
		},
		StatusCode: 400,
		Body:       string(raw),
	}, nil

}
