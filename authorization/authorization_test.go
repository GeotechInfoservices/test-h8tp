package authorization

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/xsided/h8tp/request"
	"github.com/xsided/h8tp/response"
)

func TestAuthorization(t *testing.T) {
	handler := func(ctx context.Context, req request.Request) (response.Response, error) {
		return response.OK("Test")
	}

	req := request.Request{
		Path:            "/",
		Resource:        "GET",
		HTTPMethod:      "GET",
		Body:            "",
		IsBase64Encoded: false,
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"owner_id": "",
				"scp":      "openid",
			},
		},
	}

	tt := []struct {
		Name       string
		Config     Config
		Request    request.Request
		StatusCode int
	}{
		{"No owner in token", Config{UserID: func(req request.Request) string { return "some-user-1" }}, req, 401},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			h := Authorize(handler, tc.Config)
			resp, err := h(context.Background(), tc.Request)
			if err != nil {
				t.Log("Error while trying to execute handler", err)
				t.Fail()
			}

			if resp.StatusCode != tc.StatusCode {
				t.Log("Unexpected status code, got %n, wanted %n", resp.StatusCode, tc.StatusCode)
			}
		})
	}
}
