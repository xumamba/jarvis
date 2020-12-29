/**
* @DateTime   : 2020/9/17 11:58
* @Author     : xumamba
* @Description:
**/
package server

import (
	"io"
	"math"
	"time"

	"jarvis/server/iface"
)

const abortIndex int8 = math.MaxInt8 / 2

type Context struct {
	index    int8
	handlers HandlersChain

	encoder iface.Encoder
	decoder iface.Decoder

	w io.Writer
	r io.Reader

	server *Server
	svc    *Service

	req  interface{}
	resp interface{}
	err  *Error

	container map[string]interface{}
}

func (c *Context) Next() {
	c.index++
	if c.index <= int8(len(c.handlers)) {
		c.handlers[c.index-1](c)
	}
}

func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) SetExtraValue(key string, value interface{}) {
	c.container[key] = value
}

func (c *Context) GetExtraValue(key string) (interface{}, bool) {
	value, ok := c.container[key]
	return value, ok
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *Context) Done() <-chan struct{} {
	return nil
}

func (c *Context) Err() error {
	return c.err
}

func (c *Context) Value(key interface{}) interface{} {
	return c.container[key.(string)]
}
