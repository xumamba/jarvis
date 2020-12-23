package serialize
/**
 * @DateTime   : 2020/12/23
 * @Author     : xumamba
 * @Description:
 **/

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
)


type rawJson struct {}

func (r rawJson) Marshal(src interface{}) ([]byte, error) {
	return json.Marshal(src)
}

func (r rawJson) UnMarshal(data []byte, aim interface{}) error {
	return json.Unmarshal(data, aim)
}

type jsoniterJson struct{}

func (j jsoniterJson) Marshal(src interface{}) ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(src)
}

func (j jsoniterJson) UnMarshal(data []byte, aim interface{}) error {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, aim)
}

var JSON ISerialize = jsoniterJson{}