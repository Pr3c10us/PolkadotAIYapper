package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestValidateRequest(t *testing.T) {
	validate := validator.New()

	type args struct {
		Field1 string `validate:"required"`
		Field2 int    `validate:"lte=10,gte=5"`
	}

	tests := []struct {
		name string
		args interface{}
		Err  error
	}{
		{
			name: "Valid input",
			args: &args{
				Field1: "Mr",
				Field2: 7,
			},
			Err: nil,
		},
		{
			name: "Missing required field",
			args: &args{Field2: 7},
			Err: &ValidationError{
				StatusCode: http.StatusNotAcceptable,
				Message:    "validation failed",
				ErrorMessage: []ErrorMessage{
					{"Field1", "Field1 is required"},
				},
			},
		},
		{
			name: "Field2 greater than 10",
			args: &args{Field1: "value", Field2: 11},
			Err: &ValidationError{
				StatusCode: http.StatusNotAcceptable,
				Message:    "validation failed",
				ErrorMessage: []ErrorMessage{
					{"Field2", "Field2 should be less than 10"},
				},
			},
		},
		{
			name: "Field2 less than 5",
			args: &args{Field1: "value", Field2: 4},
			Err: &ValidationError{
				StatusCode: http.StatusNotAcceptable,
				Message:    "validation failed",
				ErrorMessage: []ErrorMessage{
					{"Field2", "Field2 should be greater than 5"},
				},
			},
		},
		{
			name: "Non-validation error",
			args: errors.New("random error"),
			Err:  errors.New("random error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			var err error
			if args, ok := tt.args.(error); ok {
				err = ValidateRequest(args)
			} else {
				err = ValidateRequest(validate.Struct(tt.args))
			}

			a.Equal(tt.Err, err, "the errors should be the same")
		})
	}
}
