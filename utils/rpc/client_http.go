/**
* @Time       : 2020/9/22 9:26 下午
* @Author     : xumamba
* @Description: client_http.go
 */
package rpc

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	*http.Client
}

func NewHttpClient() IClient {
	return &HttpClient{
		&http.Client{
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 3,
		},
	}
}

func (c *HttpClient) Get(url string) (statusCode int, body []byte, err error) {
	var resp *http.Response
	resp, err = c.Client.Get(url)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return
	}
	statusCode = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func (c *HttpClient) PostText(url string, data []byte) (statusCode int, body []byte, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	req, err = http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "text/plain;charset=utf-8")
	resp, err = c.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return
	}
	statusCode = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func (c *HttpClient) PostJson(url string, data []byte) (statusCode int, body []byte, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	req, err = http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err = c.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return
	}
	statusCode = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	return
}
