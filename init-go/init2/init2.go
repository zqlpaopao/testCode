/**
 * @Author: zhangsan
 * @Description:
 * @File:  init2
 * @Version: 1.0.0
 * @Date: 2021/8/25 上午9:53
 */

package init2

import (
	"fmt"
	public_init "test/init-go/public-init"
)

func init(){
	fmt.Println("init2")
	public_init.PublicInit()
}
