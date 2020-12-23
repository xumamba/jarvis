package serialize

/**
* @DateTime   : 2020/12/23
* @Author     : xumamba
* @Description:
**/

// ISerialize 序列化
type ISerialize interface {
	Marshal(src interface{})([]byte, error)
	UnMarshal(data []byte, aim interface{}) error
}