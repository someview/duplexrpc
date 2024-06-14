package rpcinfo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type myOnewayMethod func(int) error

func TestOnewayMethodInfo(t *testing.T) {
	arg := 100
	typ := MethodType(6)
	info := NewOnewayMethodInfo(typ, func() int {
		return arg
	}, func(fn myOnewayMethod, ctx context.Context, args int) error {
		return fn(args)
	})
	assert.Equal(t, typ, info.Type())
	assert.Equal(t, arg, info.NewArg().(int))
	assert.Equal(t, nil, info.NewResult())
	assert.Equal(t, true, info.OneWay())
	info.Invoke(myOnewayMethod(func(args int) error { assert.Equal(t, arg, args); return nil }), context.Background(), info.NewArg(), info.NewResult())

}

type myMethod func(int) (int, error)

func TestMethodInfo(t *testing.T) {
	arg := 100
	resp := 200
	typ := MethodType(10)
	var fnn myMethod = func(i int) (int, error) {
		assert.Equal(t, arg, i)
		return resp, nil
	}
	info := NewMethodInfo[int, int, myMethod](typ, func() int {
		return arg
	}, func() int {
		return resp
	}, func(fn myMethod, ctx context.Context, args int, result int) error {
		re, _ := fn(args)
		assert.Equal(t, resp, re)
		return nil
	})
	assert.Equal(t, typ, info.Type())
	assert.Equal(t, arg, info.NewArg().(int))
	assert.Equal(t, resp, info.NewResult().(int))
	assert.Equal(t, false, info.OneWay())
	info.Invoke(fnn, context.Background(), info.NewArg(), info.NewResult())
}
