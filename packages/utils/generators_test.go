package utils

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_generateRandomNumber(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "success",
			args: args{
				n: 6,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newNo := GenerateRandomNumber(tt.args.n)
			t.Log(newNo)
			assert.Equalf(t, reflect.TypeOf(tt.want), reflect.TypeOf(newNo), "GenerateRandomNumber(%v)", tt.args.n)
		})
	}
}
