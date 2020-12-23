package str

/**
 * @DateTime   : 2020/12/22
 * @Author     : xumamba
 * @Description:
 **/

import (
	"fmt"
	"net"
	"sort"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// GenerateUUID 生成一个uuid
func GenerateUUID() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

// JudgeStringInSlice 判断字符串是否在目标字符串数组中
func JudgeStringInSlice(str string, arr []string) (int, bool) {
	sort.Strings(arr)
	index := sort.SearchStrings(arr, str)
	if 0 <= index && index < len(arr) && arr[index] == str {
		return index, true
	}
	return -1, false
}

// GetLocalAddr 获取本地IP地址
func GetLocalAddr() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil{
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback(){
			ip := ipNet.IP.To4()
			if ip == nil{
				continue
			}
			fmt.Println(ip.String())
			// if ip[0] != 0x0A && ip[0] != 0xC0 && ip[0] != 0xAC{
			// 	continue
			// }
			// return base64.RawStdEncoding.EncodeToString(ip), nil
		}
	}
	return "", err
}
