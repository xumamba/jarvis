/**
* @Time       : 2020/9/22 8:07 下午
* @Author     : xumamba
* @Description: http_test.go
 */
package http

import (
	"bytes"
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"

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

var (
	ln   net.Listener
	ser  *server.Server
	ts   *fasthttp.Server
	once sync.Once
)

func startServer(){
	var err error
	ser = server.DefaultServer()
	_ = ser.RegisterService("Add", func(ctx context.Context, args *Args, reply *Reply) error {
		reply.C = args.A + args.B
		return nil
	})
	ln, err = net.Listen("tcp4", ":0")
	if err != nil{
		log.Fatal(err)
	}
	ts = FastHttpServer(ser)
	go func() {
		if err := ts.Serve(ln); err != nil{
			log.Fatal(err)
		}
	}()
}

func BenchmarkHttpServer(b *testing.B) {
	once.Do(startServer)
	time.Sleep(time.Second)

	args := []byte(`{"A":10,"B":20}`)
	result := []byte("{\"C\":30}\n")

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			buf := bytebufferpool.Get()
			ser.MustHandle("Add", buf, bytes.NewBuffer(args))

			if !bytes.Equal(buf.Bytes(), result){
				b.Fatalf("expect result:%x,got:%x", result, buf.Bytes())
			}
			bytebufferpool.Put(buf)
		}
	})

}