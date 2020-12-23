package gzip
/**
 * @DateTime   : 2020/12/22
 * @Author     : xumamba
 * @Description:
 **/

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestGZip(t *testing.T) {
	srcData := []byte("测试数据压缩测试数据压缩测试数据压缩测试数据压缩测试数据压缩测试数据压缩")
	resData, err := GZip(srcData)
	fmt.Println(len(srcData), len(resData))
	assert.Nil(t, err)
	unGZip, err := UnGZip(resData)
	assert.Nil(t, err)
	assert.Equal(t, srcData, unGZip)
}

func BenchmarkGZip(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i< b.N; i++{
		GZip([]byte("第" + strconv.Itoa(i) + "次: 测试数据压缩"))
	}
}