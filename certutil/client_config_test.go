package certutil

import (
	"bytes"
	"testing"

	"go-sdk/assert"
	"go-sdk/uuid"
)

func TestNewClientConfig(t *testing.T) {
	assert := assert.New(t)

	ca, err := CreateCA()
	assert.Nil(err)

	uid := uuid.V4().String()
	client, err := CreateClient(uid, ca)
	assert.Nil(err)

	caPEM := new(bytes.Buffer)
	assert.Nil(ca.WriteCertPem(caPEM))
	clientCertPEM := new(bytes.Buffer)
	assert.Nil(client.WriteCertPem(clientCertPEM))
	clientKeyPEM := new(bytes.Buffer)
	assert.Nil(client.WriteKeyPem(clientKeyPEM))

	tlsConfig, err := NewClientConfig(
		KeyPair{Cert: clientCertPEM.String(), Key: clientKeyPEM.String()},
		[]KeyPair{{Cert: caPEM.String()}},
	)

	assert.Nil(err)
	assert.NotNil(tlsConfig)
}
