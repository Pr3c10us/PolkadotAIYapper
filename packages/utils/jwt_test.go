package utils

//
//import (
//	"github.com/Pr3c10us/gbt/internals/domain/identity"
//	"github.com/google/uuid"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestCreateTokenFromUser(t *testing.T) {
//	type args struct {
//		user   *identity.User
//		secret string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "token creation successful",
//			args: args{
//				user: &identity.User{
//					ID:               uuid.New().String(),
//					Username:         "MrMan",
//					Password:         "secretMan",
//					SecurityQuestion: "who strong",
//					SecurityAnswer:   "me",
//				},
//				secret: "secret",
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			_, err := CreateTokenFromUser(tt.args.user, tt.args.secret)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.Equal(t, err, nil)
//			}
//		})
//	}
//}
