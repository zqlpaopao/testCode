/**
 * @Author: zhangsan
 * @Description:
 * @File:  init1
 * @Version: 1.0.0
 * @Date: 2021/8/25 上午9:52
 */

package init1

import (
	"fmt"
	public_init "test/init-go/public-init"
)

func init(){
	fmt.Println("init1")
	public_init.PublicInit()
}
