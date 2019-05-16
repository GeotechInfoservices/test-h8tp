package authorization

import "github.com/aws/aws-lambda-go/events"

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

// RequireScope for the supplied context
// Makes sure the user has the necessary scopes to perform the action requested
func RequireScope(ctx events.APIGatewayProxyRequestContext, scope string) (string, error) {
	owner, ok := ctx.Authorizer["owner_id"]
	if !ok {
		return "", MissingOwnerError("invalid token provided")
	}

	return owner.(string), nil
}
