package validation

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/xsided/h8tp/request"
	"github.com/xsided/h8tp/response"
)

type Some struct {
	Other string `validate:"required" json:"other"`
}

type TestStruct struct {
	Name   string `validate:"required,email" json:"name"`
	SomeID string `validate:"required,email" json:"some_id"`
	Some   Some   `json:"some_thing"`
}

func TestValidate(t *testing.T) {
	err := ValidateStruct(TestStruct{Name: "test"})
	b, _ := json.Marshal(err)
	json := string(b)
	if json != `{"error":"validation error","validation":[{"path":"name","error":"email"},{"path":"some_id","error":"required"},{"path":"some_thing.other","error":"required"}]}` {
		t.Log()
		t.Fail()
	}
}

func TestValidateMiddleware(t *testing.T) {
	handler := func(ctx context.Context, req request.Request, test TestStruct) (response.Response, error) {
		return response.OK("Test")
	}

	Validate(handler)
}
