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
