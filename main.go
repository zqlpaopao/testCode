/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/29 下午6:42
 */

package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx,cancel := context.WithTimeout(context.Background(),2*time.Second)
	cancel= cancel
	for {
		select {
		case <- ctx.Done():
			fmt.Print("exit")
			return
		}
	}
}