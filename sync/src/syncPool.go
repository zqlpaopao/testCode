/**
 * @Author: zhangsan
 * @Description:
 * @File:  syncPool
 * @Version: 1.0.0
 * @Date: 2021/2/25 上午11:32
 */

package src

import "sync"



func Pool() {
	var sP = sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}
	sP.Get()
}