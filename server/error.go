/**
* @DateTime   : 2020/9/18 15:03
* @Author     : xumamba
* @Description:
**/
package server

type ErrorType uint64

const (
	ErrorTypeCall    ErrorType = 1 << 63
	ErrorTypePrivate ErrorType = 1 << 0
	ErrorTypePublic  ErrorType = 1 << 1
)

var _ error = &Error{}

type Error struct {
	Err  error
	Type ErrorType
	Meta interface{}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
