package certutil

import (
	"testing"

	"go-sdk/assert"
)

func TestCreateServer(t *testing.T) {
	assert := assert.New(t)

	ca, err := CreateCA()
	assert.Nil(err)

	server, err := CreateServer("warden-server", ca, OptAdditionalNames("warden-server-test"))
	assert.Nil(err)
	assert.Len(server.Certificates, 2)
	assert.Len(server.CertificateDERs, 2)
	assert.Len(server.Certificates[0].DNSNames, 2)
}
