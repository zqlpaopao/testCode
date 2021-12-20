package main

import (
	"fmt"
	"github.com/shopspring/decimal"
)

/*
	浮点型转换整型计算
 */

// 表示小数位保留8位精度
const prec = 100000000

func float2Int(f float64) int64 {
	return int64(f * prec)
}

func int2float(i int64) float64 {
	return float64(i) / prec
}

var decimalPrec = decimal.NewFromFloat(prec)

func float2Int1(f float64) int64 {
	return decimal.NewFromFloat(f).Mul(decimalPrec).IntPart()
}

func main() {
	testDecimal()
	testFloat()
	return
}

func testFloat(){
	var a, b float64 = 0.1, 0.2
	f := float2Int(a) + float2Int(b)
	fmt.Println(a+b, f, int2float(f))//0.30000000000000004 30000000 0.3
	fmt.Println(float2Int(2.3))//229999999

}

func testDecimal(){
	fmt.Println(float2Int1(2.3)) // 输出：230000000
}