package authorization

import (
	"context"
	"testing"

	"dev.azure.com/ManyDigital/MDLiveQuiz/_git/livequiz-h8tp/request"
	"dev.azure.com/ManyDigital/MDLiveQuiz/_git/livequiz-h8tp/response"
	"github.com/aws/aws-lambda-go/events"
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
				// "owner_id": "",
				"role": "user",
				// "scope":    "openid,auth",

				"owner_id": "stephans-owner",
				// "role":     "user",
				"scope": "openid,results,events,games,fanshout,users,sponsors",
			},
		},
	}

	tt := []struct {
		Name       string
		Config     Config
		Request    request.Request
		StatusCode int
	}{
		{"No owner in token", Config{}, req, 500},
		{"Not required scope", Config{RequiredScope: "non-existant", Role: Administartor}, req, 401},
		{"required scope", Config{RequiredScope: "auth"}, req, 200},
		{"required role", Config{RequiredScope: "xyz", Role: Tester}, req, 200},
		{"wrong role", Config{RequiredScope: "xyz", Role: Administartor}, req, 401},
		{"no role required but has role", Config{RequiredScope: "xasdyz"}, req, 200},
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
				t.Logf("%s, Unexpected status code, got %d, wanted %d", tc.Name, resp.StatusCode, tc.StatusCode)
				t.Fail()
			}
		})
	}
}
