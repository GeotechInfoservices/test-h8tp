package authorization

import (
	"context"
	"fmt"
	"strings"

	"dev.azure.com/ManyDigital/MDLiveQuiz/_git/livequiz-h8tp/request"
	"dev.azure.com/ManyDigital/MDLiveQuiz/_git/livequiz-h8tp/response"
	"github.com/aws/aws-lambda-go/events"
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
	Role          role
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

// Contains tests if a string exists in a scope
func Contains(required string, scopes []string) bool {
	for _, scope := range scopes {
		if scope == required {
			return true
		}
	}

	return false
}

type role int

// Roles
const (
	God           role = 6
	Administartor role = 5
	Publisher     role = 4
	Operator      role = 3
	Launcher      role = 2
	Tester        role = 1
)

// Role tests if a role is present or higher in hierachy
func Role(required role, r string) bool {
	roles := map[string]role{
		"tester":        Tester,
		"launcher":      Launcher,
		"operator":      Operator,
		"publisher":     Publisher,
		"administrator": Administartor,
		"god":           God,
	}

	actual, ok := roles[r]
	if ok {
		if actual >= required {
			return true
		}
	}

	return false
}

// Authorize http request
// Checks the request context for an owner id and performs checks based on the given config
func Authorize(h func(context.Context, request.Request) (response.Response, error), c Config) func(context.Context, request.Request) (response.Response, error) {
	return func(ctx context.Context, req request.Request) (response.Response, error) {

		if c.RequiredScope == "" {
			fmt.Println("Misconfigured authorization")
			return response.InternalServerError()
		}

		_, ok := req.RequestContext.Authorizer["owner_id"]
		if !ok {
			return response.Unauthorized("invalid token provided")
		}

		scopes := strings.Split(req.RequestContext.Authorizer["scope"].(string), ",")
		allow := Contains(c.RequiredScope, scopes)

		if allow == true {
			return h(ctx, req)
		}

		r, ok := req.RequestContext.Authorizer["role"]
		if ok {
			allow = Role(c.Role, r.(string))
			if allow == true {
				return h(ctx, req)
			}
		}

		return response.Unauthorized("unauthorized")
	}
}
