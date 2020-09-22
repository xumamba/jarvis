/**
* @Time       : 2020/9/22 9:21 下午
* @Author     : xumamba
* @Description: client.go
 */
package rpc

type IClient interface {
	Get(url string) (statusCode int, body []byte, err error)
	PostText(url string, data []byte) (statusCode int, body []byte, err error)
	PostJson(url string, data []byte) (statusCode int, body []byte, err error)
}
