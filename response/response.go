package response

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/xsided/h8tp/validation"
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

// ValidationError for handling validation errors
func ValidationError(resp *validation.Error) (Response, error) {
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
		StatusCode: 400,
		Body:       string(b),
	}, nil
}
