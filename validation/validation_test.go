package validation

import (
	"encoding/json"
	"testing"
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
