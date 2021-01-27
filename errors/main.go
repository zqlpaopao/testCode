/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/26 下午5:51
 */

package main

import (
	"fmt"
	"errors"
)

func main() {
	e1 := errors.New("err1")
	e2 := fmt.Errorf("Wrap了一个错误err2 - %w", e1)
	e3 := fmt.Errorf("Wrap了一个错误err3 - %w", e2)
	fmt.Println(e2)
	fmt.Println(e3)
	fmt.Println(errors.Unwrap(e3))
	fmt.Println(errors.Unwrap(errors.Unwrap(e3)))

	e4 := errors.New("jkhg")
	fmt.Println(errors.As(e3,&e4))
	if errors.As(e3,&e4){
		fmt.Println("as is have")
	}

	if errors.Is(e3,e1){
		fmt.Println("is have")
	}
}

