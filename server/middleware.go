/**
* @DateTime   : 2020/9/17 13:38
* @Author     : xumamba
* @Description:
**/
package server

import (
	"fmt"
	"log"
	"reflect"

	"github.com/facebookgo/stack"
	jsoniter "github.com/json-iterator/go"
)

type HandlerFunc func(c *Context)

type HandlersChain []HandlerFunc

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func JsonHandler() HandlerFunc {
	return func(c *Context) {
		c.encoder = jsoniter.NewEncoder(c.w)
		c.decoder = jsoniter.NewDecoder(c.r)
		c.Next()
	}
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				s := stack.Caller(3)
				if IsDebugMode() {
					fmt.Printf("[Recovery] panic recovered:\nArgs:%v\nError:%s\n%s\n", c.req, err, s)
					c.err = &Error{Err: err.(error), Type: ErrorTypePublic, Meta: s.String()}
				} else {
					c.err = &Error{Err: err.(error), Type: ErrorTypePublic}
				}

				c.Abort()
				if err = c.encoder.Encode(c.err); err != nil {
					log.Println(err)
				}
			}
		}()
		c.Next()
	}
}

func HandleRequest() HandlerFunc {
	return func(c *Context) {
		c.req = New(c.svc.reqType)
		if !c.svc.IsGRPC {
			c.resp = New(c.svc.respType)
		}
		if err := c.decoder.Decode(c.req); err != nil {
			panic(err)
		}

		c.Next()

		if err := c.encoder.Encode(c.resp); err != nil {
			log.Println(err)
		}

	}
}

func Call() HandlerFunc {
	return func(c *Context) {
		var err error
		if c.svc.IsGRPC {
			err = c.svc.grpcCall(c)
		} else {
			err = c.svc.call(c, c.req, c.resp)
		}
		if err != nil {
			panic(fmt.Errorf("%s\n%s", c.svc.name, err))
		}
		c.Next()
	}
}

func DefaultServer() *Server {
	return NewServer().Use(Recovery(), JsonHandler(), HandleRequest(), Call())
}

func New(t reflect.Type) interface{} {
	var req reflect.Value

	if t.Kind() == reflect.Ptr {
		req = reflect.New(t.Elem())
	} else {
		req = reflect.New(t)
	}

	return req.Interface()
}
