package certutil

import (
	"testing"

	"go-sdk/assert"
	"go-sdk/uuid"
)

func TestCreateClient(t *testing.T) {
	assert := assert.New(t)

	ca, err := CreateCA()
	assert.Nil(err)

	uid := uuid.V4().String()
	client, err := CreateClient(uid, ca)
	assert.Nil(err)
	assert.Len(client.Certificates, 2)
	assert.Len(client.CertificateDERs, 2)
}
