package str

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
* @DateTime   : 2020/12/22
* @Author     : xumamba
* @Description:
**/

func TestStr(t *testing.T)  {
	uuid := GenerateUUID()
	fmt.Println(uuid)

	index, ok := JudgeStringInSlice("aa", []string{"bb", "cc"})
	assert.Equal(t, index, -1)
	assert.False(t, ok)
	index, ok = JudgeStringInSlice("cc", []string{"bb", "cc"})
	assert.Equal(t, index, 1)
	assert.True(t, ok)

	localAddr, err := GetLocalAddr()
	assert.Nil(t, err)
	fmt.Println(localAddr)
}