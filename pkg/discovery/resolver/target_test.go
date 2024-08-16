package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseServiceAddr(t *testing.T) {
	res, err := parseServiceAddr("local:///jd-gopush:5002")
	assert.Nil(t, err)
	assert.Equal(t, res.serviceName, "jd-gopush")
	assert.Equal(t, res.port, "5002")
	// assert.Equal(t, res.serviceNamespace, defaultNameSpace)
}

func Test_parseServiceName(t *testing.T) {
	service, port, err := parseServiceName("local:///jd-gopush:5002")
	assert.Nil(t, err)
	assert.Equal(t, service, "jd-gopush")
	assert.Equal(t, port, "5002")
	service, port, err = parseServiceName("jd-gopush:5002")
	assert.Nil(t, err)
	assert.Equal(t, service, "jd-gopush")
	assert.Equal(t, port, "5002")
	service, port, err = parseServiceName("ws://jd-gopush:5002")
	assert.NotNil(t, err)
}
