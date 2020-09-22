/**
* @Time       : 2020/9/22 8:07 下午
* @Author     : xumamba
* @Description: http_test.go
 */
package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"jarvis/server"
	"jarvis/utils/rpc"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

func TestServer(t *testing.T) {
	s := server.DefaultServer()
	err := s.RegisterService("Add", func(ctx context.Context, args *Args, reply *Reply) error {
		reply.C = args.A + args.B
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	s.MustHandle("Add", os.Stdout, bytes.NewBufferString(`{"A": 1,"B": 1}`))
}

func TestFastHttpServer(t *testing.T) {
	s := server.DefaultServer()
	err := s.RegisterService("Add", func(ctx context.Context, args *Args, reply *Reply) error {
		reply.C = args.A + args.B
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(Handler(s))
	defer ts.Close()

	httpClient := rpc.NewHttpClient()
	statusCode, body, err := httpClient.PostJson(ts.URL+"/Add", []byte(`{"A": 1,"B": 1}`))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "{\"C\":2}\n", string(body))
}
