package authorization

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
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

// UserEdit config. Wether or not a user should be allowed to edit his own data
type UserEdit bool

// Permissions for user edit.
const (
	Allowed UserEdit = true  // Allow user to edit his own data
	Denied           = false // Deny user to edit his own data
)

// Config for authorization
type Config struct {
	User   UserEdit
	Scope  string
	UserID func(req request.Request) string
}

// RequireScope for the supplied context
// Makes sure the user has the necessary scopes to perform the action requested
func RequireScope(ctx events.APIGatewayProxyRequestContext, scope string) (string, error) {
	owner, ok := ctx.Authorizer["owner_id"]
	if !ok {
		return "", MissingOwnerError("invalid token provided")
	}

	return owner.(string), nil
}

// CurrentUser for the supplied context
// Returns the user currently authenticated
func CurrentUser(ctx events.APIGatewayProxyRequestContext) string {
	logrus.Infof("Context: %+v", ctx)
	userID := ctx.Authorizer["sub"]

	return userID.(string)
}

// Handler for http requests
type Handler func(request.Request) (events.APIGatewayProxyResponse, error)

// AuthorizeWithOwner authorizes an http request.
// Checks the request context for an owner id and performs checks based on the given config
func AuthorizeWithOwner(h Handler, c Config) Handler {
	return func(req request.Request) (events.APIGatewayProxyResponse, error) {
		userID := c.UserID(req)

		if userID != CurrentUser(req.RequestContext) {
			_, err := RequireScope(req.RequestContext, c.Scope)
			if err != nil {
				logrus.Errorf("%+v", err)
				return response.Unauthorized("not allowed")
			}
		}

		_, err := RequireScope(req.RequestContext, "")
		if err != nil {
			logrus.Errorf("%+v", err)
			return response.Unauthorized("not allowed")
		}

		return h(req)
	}
}
