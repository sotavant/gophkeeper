package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const cryptKeysPath = "/Users/fanishadeev/GolandProjects/gophkeeper/internal/crypto"

func TestCipher_Encrypt(t *testing.T) {
	c, err := NewCipher(cryptKeysPath)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	testMsg := "test message"

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got, decoded []byte
			got, err = c.Encrypt([]byte(testMsg))
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			decoded, err = c.Decrypt(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, string(decoded), testMsg)
		})
	}
}

func TestCipher_GetGRPCTransportCreds(t *testing.T) {
	tests := []struct {
		name         string
		cert         string
		wantProtocol string
	}{
		{
			name:         "with certificate",
			cert:         certPath,
			wantProtocol: "tls",
		},
		{
			name:         "without certificate",
			cert:         "",
			wantProtocol: "insecure",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cipher{
				certPath: tt.cert,
			}

			assert.Equal(t, tt.wantProtocol, c.GetClientGRPCTransportCreds().Info().SecurityProtocol)
		})
	}
}

func TestCipher_GetServerGRPCTransportCreds(t *testing.T) {
	tests := []struct {
		name         string
		cert         string
		wantProtocol string
	}{
		{
			name:         "with certificate",
			cert:         certPath,
			wantProtocol: "tls",
		},
		{
			name:         "without certificate",
			cert:         "",
			wantProtocol: "insecure",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cipher{
				certPath: tt.cert,
			}

			assert.Equal(t, tt.wantProtocol, c.GetClientGRPCTransportCreds().Info().SecurityProtocol)
		})
	}
}
