package certutil

import (
	"testing"

	"go-sdk/assert"
)

func TestCreateCA(t *testing.T) {
	assert := assert.New(t)

	ca, err := CreateCA()
	assert.Nil(err)
	assert.NotNil(ca.PrivateKey)
	assert.NotNil(ca.PublicKey)
	assert.Len(ca.Certificates, 1)
	assert.Len(ca.CertificateDERs, 1)
}
