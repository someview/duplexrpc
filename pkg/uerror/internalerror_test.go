package uerror

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTransError(t *testing.T) {
	err1 := fmt.Errorf("123：%w", NewTransError(111, "456"))
	assert.True(t, IsTransError(err1))
	var err2 error = &BasicError{
		errorCategory: ErrTrans,
		code:          0,
		message:       "",
	}
	assert.True(t, IsTransError(err2))
	var err3 = fmt.Errorf("test error")
	assert.False(t, IsTransError(err3))
}
