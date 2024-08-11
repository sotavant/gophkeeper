package data

import (
	"gophkeeper/client/domain"
	"gophkeeper/internal/client"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_hashData(t *testing.T) {
	client.AppInstance = &client.App{}
	client.AppInstance.SetStorageKey("login", "pass")

	type args struct {
		data domain.Data
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test hash data",
			args: args{
				data: domain.Data{
					Name:  "name",
					Login: "login",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encryptData(tt.args.data)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.data.Name, got.Name)
			assert.NotEqual(t, tt.args.data.Login, got.Login)
			assert.Equal(t, "", got.Pass)
		})
	}
}
