package response

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
	"github.com/aws/aws-lambda-go/events"
)

// InvalidRequest response for API. This automatically wraps the error message in the correct format.
func InvalidRequest(msg string) (events.APIGatewayProxyResponse, error) {

	resp := map[string]interface{}{
		"error": msg,
	}

	// ignore error, since we know this will always pass given a simple string.
	b, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Origin":      "*",
		},
		StatusCode: 400,
		Body:       string(b),
	}, nil
}

// Unauthorized response for API. This automatically wraps the error message in the correct format.
func Unauthorized(msg string) (events.APIGatewayProxyResponse, error) {

	resp := map[string]interface{}{
		"error": msg,
	}

	// ignore error, since we know this will always pass given a simple string.
	b, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Origin":      "*",
		},
		StatusCode: 401,
		Body:       string(b),
	}, nil
}

// OK 200 response from api. This automatically marshals a struct and converts it to json
func OK(body interface{}) (events.APIGatewayProxyResponse, error) {
	resp := map[string]interface{}{
		"data": body,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		msg := map[string]string{"error": "error while retrieving items"}
		b, _ := json.Marshal(msg)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       string(b),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Origin":      "*",
		},
		StatusCode: 200,
		Body:       string(b),
	}, nil
}

type Error struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationError struct {
	Error  string  `json:"error"`
	Errors []Error `json:"validation"`
}

// BadInput implies an error in the input, according to the entity validation rules.
func BadInput(errors error) (events.APIGatewayProxyResponse, error) {
	out := ValidationError{
		Error: "invalid input",
	}
	switch err := errors.(type) {
	case govalidator.Errors:
		for _, e := range err.Errors() {
			et := e.(govalidator.Error)
			out.Errors = append(out.Errors, Error{Error: et.Error(), Field: et.Name})
		}
	}

	raw, err := json.Marshal(out)
	if err != nil {
		// TODO: Do some checking
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Origin":      "*",
		},
		StatusCode: 400,
		Body:       string(raw),
	}, nil

}
