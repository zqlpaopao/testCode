/**
 * @Author: zhangsan
 * @Description:
 * @File:  gos_test.go
 * @Version: 1.0.0
 * @Date: 2021/3/25 下午2:12
 */

package src

import (
	"testing"
	"go.uber.org/goleak"
)

func TestGo1(t *testing.T) {
	defer goleak.VerifyNone(t)
	Go1()
}
