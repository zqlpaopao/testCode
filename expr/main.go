package main

import (
	"fmt"

	"github.com/expr-lang/expr"
)

func main() {
	// 环境变量定义
	env := map[string]interface{}{
		"1": false,
		"2": true,
		"3": false,
	}

	// 直接编写逻辑表达式
	code := `1 && (2 || 3)`

	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Printf("结果: %v\n", output) // 结果: false
}
