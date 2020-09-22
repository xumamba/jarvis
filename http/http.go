/**
* @DateTime   : 2020/9/22 9:22
* @Author     : xumamba
* @Description:
**/
package http

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/facebookgo/stack"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"

	"jarvis/server"
)

func timeFormat(t time.Time) string {
	timeStr := t.Format("2020/09/22 09:22:38")
	return timeStr
}

func FastHttpServer(s *server.Server) *fasthttp.Server {
	reqHandler := func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				s := stack.Caller(3)
				if server.IsDebugMode() {
					log.Printf("[Recovery] %s panic recovered:\n%s\n%s\n%s", timeFormat(time.Now()), ctx.Request.String(), err, s)
				} else {
					log.Printf("[Recovery] %s panic recovered:\n%s\n%s", timeFormat(time.Now()), err, s)
				}

				ctx.SetStatusCode(http.StatusInternalServerError)
			}
		}()
		svcName := string(ctx.Path())[1:]
		byteBuffer := bytebufferpool.Get()
		s.MustHandle(svcName, byteBuffer, bytes.NewReader(ctx.PostBody()))
		if _, err := ctx.Write(byteBuffer.Bytes()); err != nil {
			log.Printf("write response error:%s", err)
		}
		bytebufferpool.Put(byteBuffer)
		ctx.SetStatusCode(http.StatusOK)
	}
	return &fasthttp.Server{Handler: reqHandler}
}

func Handler(s *server.Server) http.Handler {
	mux := http.NewServeMux()
	s.ForEachService(func(name string, svc *server.Service) {
		mux.HandleFunc("/"+name, func(responseWriter http.ResponseWriter, request *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					s := stack.Caller(3)
					if server.IsDebugMode() {
						log.Printf("[Recovery] %s panic recovered:\n%s\n%s\n%s", timeFormat(time.Now()), ctx.Request.String(), err, s)
					} else {
						log.Printf("[Recovery] %s panic recovered:\n%s\n%s", timeFormat(time.Now()), err, s)
					}

					responseWriter.WriteHeader(http.StatusInternalServerError)
				}
			}()

			s.MustHandle(name, responseWriter, request.Body)
		})
	})
	return mux
}
