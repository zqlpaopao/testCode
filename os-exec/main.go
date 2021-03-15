/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/3/15 下午2:05
 */

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main(){
	//每个shell 要分割输入
	cmd := exec.Command("git","branch")
	c,err := cmd.Output()
	fmt.Println(strings.Trim(string(c),"\r\n"))
	fmt.Println(err)
}
