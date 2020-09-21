/**
* @DateTime   : 2020/9/18 15:13
* @Author     : xumamba
* @Description:
**/
package server

type Encoder interface {
	Encode(val interface{}) error
}

type Decoder interface {
	Decode(obj interface{}) error
}