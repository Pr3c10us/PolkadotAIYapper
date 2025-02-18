package command

import (
	"github.com/Pr3c10us/boilerplate/internals/domains/embedding"
	"github.com/Pr3c10us/boilerplate/internals/domains/llm"
	"github.com/Pr3c10us/boilerplate/internals/domains/xdotcom"
	"reflect"
	"testing"
)

func TestTweet_convertToArray(t *testing.T) {
	type fields struct {
		llm       llm.Repository
		embedding embedding.Repository
		xdotcom   xdotcom.Repository
	}
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Tweet{
				llm:       tt.fields.llm,
				embedding: tt.fields.embedding,
				xdotcom:   tt.fields.xdotcom,
			}
			got, err := service.convertToArray(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToArray() got = %v, want %v", got, tt.want)
			}
		})
	}
}
