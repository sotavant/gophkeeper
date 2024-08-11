package crypto

import (
	"crypto/sha1"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/pbkdf2"
)

func TestEncrypt(t *testing.T) {
	testText := "hello world"

	type args struct {
		key  []byte
		data []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test success",
			args: args{
				key:  []byte("some_key"),
				data: []byte(testText),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := pbkdf2.Key(tt.args.key, []byte("some salt"), 4096, 32, sha1.New)

			got, err := Encrypt(key, tt.args.data)
			assert.NoError(t, err)

			decrypted, err := Decrypt(key, got)
			assert.NoError(t, err)

			assert.Equal(t, testText, decrypted)
		})
	}
}
