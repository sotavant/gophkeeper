package auth

import (
	"fmt"
	"gophkeeper/internal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildJWTString(t *testing.T) {
	got, err := BuildJWTString(3333)
	assert.NoError(t, err)
	assert.NotEmpty(t, got)
}

func TestGetUserID(t *testing.T) {
	var got uint64
	var id uint64 = 3333
	internal.InitLogger()

	token, err := BuildJWTString(id)
	assert.NoError(t, err)

	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				tokenString: token,
			},
			want:    id,
			wantErr: assert.NoError,
		},
		{
			name: "bad token string",
			args: args{
				tokenString: "sdfffff",
			},
			want:    0,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err = GetUserID(tt.args.tokenString)
			if !tt.wantErr(t, err, fmt.Sprintf("GetUserID(%v)", tt.args.tokenString)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetUserID(%v)", tt.args.tokenString)
		})
	}
}
