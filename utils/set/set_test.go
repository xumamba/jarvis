package set

/**
 * @DateTime   : 2020/12/23
 * @Author     : xumamba
 * @Description:
 **/

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	s := NewSet()
	s.Add([]interface{}{1, 2, 3, 4, 5, 3, 2}...)
	fmt.Println(s.List())
	assert.Equal(t, s.Len(), 5)
	assert.True(t, s.Exists(1))
	s.Remove(2)
	assert.False(t, s.Exists(2))
	s.Clear()
	fmt.Print(s.List())
}
