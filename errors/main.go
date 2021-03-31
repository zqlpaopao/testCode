package main

import (
	"github.com/pkg/errors"
	"fmt"
	"os"
	"test/errors/src"

	//"errors"
)

func f() error {
	return errors.New("is f")
}

func errorWith()error{
	return  f()
	//return errors.WithStack(errors.New("ggg"))
}

func main() {
	// 详细输出（打印出调用堆栈）
	//fmt.Printf("%+v\n", f())
	//aa := fmt.Sprintf("%+v\n", errorWith())
	//fmt.Println(aa)

	var myError = &src.MyError{}
	var err error = src.Foo()
	for err != nil {
		containNotExist := errors.Is(err, os.ErrNotExist)
		isMyError := errors.As(err, &myError)
		fmt.Println(containNotExist, isMyError)
		fmt.Printf("%+v\n", err)
		err = errors.Unwrap(err)
	}
	fmt.Println("=========================")

	err = src.Foo().Unwrap()
	fmt.Printf("%+v\n", err)
	err = errors.Unwrap(err)
	fmt.Printf("%+v\n", err)
}