/**
* @DateTime   : 2020/9/18 15:33
* @Author     : xumamba
* @Description:
**/
package server

import (
	"context"
	"fmt"
	"reflect"
)

type Service struct {
	name     string
	fun      reflect.Value
	reqType  reflect.Type
	respType reflect.Type
	IsGRPC   bool
}

func (svc *Service) call(ctx context.Context, req, resp interface{}) (err error) {
	returnValues := svc.fun.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req), reflect.ValueOf(resp)})

	if errInter := returnValues[0].Interface(); errInter != nil {
		err = errInter.(error)
	}
	return
}

func (svc *Service) grpcCall(ctx context.Context) (err error) {
	c, ok := ctx.(*Context)
	if !ok {
		panic("unknown Context")
	}

	returnValues := svc.fun.Call([]reflect.Value{reflect.ValueOf(c), reflect.ValueOf(c.req)})
	callErr := returnValues[1].Interface()
	if callErr != nil {
		err = callErr.(error)
		return
	}

	c.resp = returnValues[0].Interface()
	if c.resp == nil {
		err = fmt.Errorf("response is a nil value")
	}
	return
}
