package authorization

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/xsided/h8tp/request"
	"github.com/xsided/h8tp/response"
)

// MissingOwnerError indicates a context was provided without an owner id present
func MissingOwnerError(msg string) *MissingOwner {
	return &MissingOwner{
		msg: msg,
	}
}

// MissingOwner struct for the MissingOwnerError
type MissingOwner struct {
	msg string
}

func (e *MissingOwner) Error() string {
	return e.msg
}

// Config for authorization
type Config struct {
	RequiredScope string
	OverrideScope string
	UserID        func(req request.Request) string
}

// GetOwner for the supplied context
// Makes sure the user has the necessary scopes to perform the action requested
func GetOwner(req request.Request) string {
	owner := req.RequestContext.Authorizer["owner_id"]

	return owner.(string)
}

// CurrentUser for the supplied context
// Returns the user currently authenticated
func CurrentUser(ctx events.APIGatewayProxyRequestContext) string {
	userID := ctx.Authorizer["principalId"]

	return userID.(string)
}

// Authorize http request
// Checks the request context for an owner id and performs checks based on the given config
func Authorize(h func(context.Context, request.Request) (response.Response, error), c Config) func(context.Context, request.Request) (response.Response, error) {
	return func(ctx context.Context, req request.Request) (response.Response, error) {

		_, ok := req.RequestContext.Authorizer["owner_id"]
		if !ok {
			return response.Unauthorized("invalid token provided")
		}

		// req.RequestContext.Authorizer["role"] = "user"

		// TODO: Handle roles
		// TODO: Handle correct scope in token, or deny access

		return h(ctx, req)
	}
}
