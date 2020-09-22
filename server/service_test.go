/**
* @Time       : 2020/9/21 8:21 下午
* @Author     : xumamba
* @Description: service_test.go
 */
package server

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct{}

func TestService(t *testing.T) {
	ast := assert.New(t)
	server := NewServer()

	err := server.RegisterService("testSvcName", func() {})
	ast.EqualError(err, "register function parameters must be:func Func(ctx context.Context, argv *Arg, reply *Reply) error")

	err = server.RegisterService("testSvcName", func(ctx context.Context, req testStruct) error { return nil })
	ast.EqualError(err, "register function parameters must be:func Func(ctx context.Context, argv *Arg, reply *Reply) error")

	err = server.RegisterService("testSvcName", func(ctx context.Context, req, resp testStruct) error { return nil })
	ast.EqualError(err, "register function arg type must be ptr and export:server.testStruct")

	err = server.RegisterService("testSvcName", func(ctx context.Context, req, resp *testStruct) error { return nil })
	ast.EqualError(err, "register function arg type must be ptr and export:*server.testStruct")

	err = server.RegisterService("testSvcName", func(ctx context.Context, req, resp *string) testStruct { return testStruct{} })
	ast.EqualError(err, "register function return type must be error:server.testStruct")

	err = server.RegisterService("testSvcName", func(ctx context.Context, req, resp *string) error {
		*resp = strings.ToUpper(*req)
		return nil
	})
	ast.Nil(err)

	svc, ok := server.serviceMap.Load("testSvcName")
	ast.True(ok)
	var req, resp string
	req = "jarvis"
	err = svc.(*Service).call(context.Background(), &req, &resp)
	ast.Equal(resp, "JARVIS")
}
