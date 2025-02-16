package utils

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestIsValidPassword(t *testing.T) {
	pw := "secretMan"
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	if err != nil {
		t.Error(err)
	}
	type args struct {
		hashedPw string
		pw       string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "correct password",
			args: args{
				hashedPw: string(hashedPw),
				pw:       "secretMan",
			},
			want: true,
		},
		{
			name: "incorrect password",
			args: args{
				hashedPw: string(hashedPw),
				pw:       "wrongPassword",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPassword(tt.args.hashedPw, tt.args.pw); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
