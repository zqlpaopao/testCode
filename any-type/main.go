package main

import (
	"fmt"
	"strconv"
)

//https://mp.weixin.qq.com/s/oA5MVUg1dOS5iD_hE5V1Fw

func main() {
	Print("字符串")
	Print2(2)
	Print3(2.00)

	//**********************************泛型方法的调用参数限制****************************//
	//方法参数限制，如果制定必须符合泛型的规定，可以制定，或者[]在方法前标识
	Sub(32, 12)
	Sub(int64(32), int64(12))

	//调用泛型方法,可以
	PrintAnyFunc[int]([]int{1, 2, 3})

	//内置可比类型
	//参数==的!=值 comparable
	i := Index([]int{1, 2, 3}, 2)
	fmt.Println("---i", i)
	i1 := Index([]string{"111", "222"}, "111")
	fmt.Println("字符串-i1", i1)

	//Cannot use map[string]string as the type comparable
	//i12 := Index([]map[string]string{{
	//	"aa": "aa",
	//	"bb": "cc",
	//}}, map[string]string{
	//	"aa": "aa",
	//})
	//fmt.Println("切片map-i1", i12)

	//3、类型推断
	fmt.Println("3、类型推断")
	E()

	//4、闭包函数
	fmt.Println("4、闭包函数")

	bar := func() string {
		return "bar"
	}
	foobar := foo(bar)
	fmt.Println(foobar())

	//5、闭包函数--filter
	fmt.Println("5、闭包函数--filter")
	GetFilterSlice()

	//5、指针类型--filter
	fmt.Println("5、指针类型--filter")
	F()
}

//第一种

func Print[T any](t T) {
	fmt.Println("第一种", t)
}

//第二种

type T any

func Print2(t T) {
	fmt.Println("第二种", t)
}

//第三种

func Print3(t any) {
	fmt.Println("第三种", t)

}

// 限定为 Types

func Sub[T Types](t1, t2 T) T {
	return t1 - t2
}

//********************************泛型方法********************

func PrintAnyFunc[T any](s []T) {
	for _, v := range s {
		fmt.Println(v)
	}
}

//Cannot use 'v' (type T) as the type string
//没有any 具体的类型，不能
//
//func Stringify[T any](s []T) (ret []string) {
//	for _, v := range s {
//		ret = append(ret, v) // 编译错误
//	}
//	return ret
//}

//********************************接口类型********************
//因为在ConcatTo中 调用了v.string，所以要求必须实现string，不然会报错

type Stringer interface {
	String() string
}

type Plusser interface {
	Plus(string) string
}

func ConcatTo[S Stringer, P Plusser](s []S, p []P) []string {
	r := make([]string, len(s))
	for i, v := range s {
		r[i] = p[i].Plus(v.String())
	}
	return r
}

//********************************类型集合********************

// 联合约束 写成一系列由竖线 ( |) 分隔的约束元素
//并集元素的类型集是序列中每个元素的类型集的并集。联合中列出的元素必须全部不同 底层类型相同

type Types interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 |
	~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64
}

//********************************约定的可比类型comparable********************
//Go1.18 中内置了一个类型约束 comparable约束，comparable约束的类型集是所有可比较类型的集合。
//这允许使用该类型参数==的!=值。
//booleans, numbers, strings, pointers, channels, arrays of comparable types,
// structs whose fields are all comparable types

type Dictionay[K comparable, V any] map[K]V

func Index[T comparable](s []T, x T) bool {

	//********************** 可以直接定义个map类型，然后都可以用户这个
	dict := Dictionay[string, int]{
		"string": 1,
	}

	d := Dictionay[string, string]{
		"string": "string",
	}
	fmt.Printf("dict: %#v \n", dict)
	fmt.Printf("dict: %#v \n", d)
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

//在接口中使用
//key类型必须定义比较操作符==和!=；
//key类型必须不能是function、map或slice（没有定义比较操作符）；
//对于interface类型，其动态类型必须定义比较操作符；
//不满足上述约束，则会导致运行时异常（run-time panic）。
//

type ComparableHasher interface {
	comparable
	Hash() uintptr
}

//interface includes constraint elements, can only be used in type parameters
//type Ts struct {
//	c comparable
//
//}

//********************************类型推断********************

func Map[F, T any](s []F, f func(F) T) {
	for _, v := range s {
		f(v)
	}
}

func E() {
	var s []int
	s = []int{1, 3, 5, 6}
	f := func(i int) int64 {
		fmt.Println("----i", i)
		return int64(i)
	}

	// 标注两个类型
	Map[int, int64](s, f)
	// 只指定第一个类型参数
	Map[int](s, f)
	// 不指定任何类型参数，并让两者都被推断。
	Map(s, f)
}

// ********************************泛型函数式应用********************
func foo(bar func() string) func() string {
	return func() string {
		return "foo" + " " + bar()
	}
}

// ********************************闭包-filter函数********************

type Model interface {
	~int | ~int64 | ~int32
}

type TModel any

//type Model interface {
//	comparable
//}

type Func func(model int) bool

func Filter[T any](f func(T) bool, src []T) []T {
	var dest []T
	for _, v := range src {
		if f(v) {
			dest = append(dest, v)
		}
	}
	return dest
}

func GetFilterSlice() {
	src := []int{2, 4, 6, 8}

	ints := Filter(func(t1 int) bool {
		if t1 > 0 {
			return true
		}
		return false
	}, src)

	fmt.Println("ints-------", ints)

	src1 := []string{"aa", "bb", "cc"}

	strs := Filter(func(t1 string) bool {
		if t1 != "" {
			return true
		}
		return false
	}, src1)

	fmt.Println("strs", strs)
}

//********************************* ge 指针 类型 ************************

//代理者模式的时候可以使用

type Setter interface {
	Set(string)
}

func FromStrings[T Setter](s []string) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[i].Set(v)
	}
	return result
}

type SetTable int

func (s *SetTable) Set(string2 string) {
	i, _ := strconv.Atoi(string2)
	s = new(SetTable)
	*s = SetTable(i)
}

//指针类型*Settable实现了约束，但代码确实想使用非指针类型Settable。
//我们需要的是一种编写方法FromStrings，它可以将类型Settable作为参数但调用指针方法。
//重复一遍，我们不能使用Settable，因为它没有Set方法，我们不能使用*Settable，
//因为不能创建 type 的切片Settable。

type Setter2[B any] interface {
	Set(string)
	//非接口类型约束元素
	*B // non-interface type constraint element
}

func FromStrings2[T any, PT Setter2[T]](s []string) []T {
	result := make([]T, len(s))
	for i, v := range s {
		p := PT(&result[i])
		p.Set(v)
	}
	return result
}

func F() {
	//Cannot use SetTable as the type Setter Type
	//does not implement 'Setter' as the 'Set' method has a pointer receiver
	//函数F试图用转换返回类型为Settable，
	//但Settable没有Set方法。有Set方法的类型是*Settable，那在调用时将返回类型改变为*Settable。
	//nums := FromStrings[SetTable]([]string{"1", "2"})

	nums := FromStrings[*SetTable]([]string{"1", "2"})

	fmt.Println("nums", nums)

	nums1 := FromStrings2[SetTable, *SetTable]([]string{"1", "2"})
	fmt.Println("nums1", nums1)
}

//********************************* 泛型  0值 ************************
//Tips:
//
//Go中现有泛型设计对于类型参数的零值并不好表达，Go官方目前没有更好的办法，但是提供了一些目前可行的一些方案：
//
//对于目前泛型的设计：
//
//可用 var zero T，但是这里需要额外去声明下。
//使用*new(T)。
//对于返回结果可命名结果参数，并使用裸return返回零值。
//扩展设计：
//
//设计以允许nil用作任何泛型类型的零值（但请参阅issue 22729[13]）。
//设计以允许使用T{}（其中T是类型参数）来指示类型的零值。
//更改语言以允许return ...返回结果类型的零值，如issue 21182[14]中所建议的那样。
//但目前来说一般使用 var zero T 的方式。

//********************************* 泛型  0值 ************************

type Queue[t any] struct {
	data chan t
}

//构建新队列

func NewQueue[T any](size int) Queue[T] {
	return Queue[T]{
		data: make(chan T, size),
	}
}

// Push 压入数据
func (q Queue[T]) Push(val T) {
	q.data <- val
}

// Pop 弹出数据 ,如果没有数据会被阻塞
func (q Queue[T]) Pop() T {
	d := <-q.data
	return d
}

func (q Queue[T]) TryPop() (T, bool) {
	select {
	case val := <-q.data:
		if val == nil {

		}
		return val, true
	default:
		// 编译报错
		//Cannot use 'nil' as the type T
		//****** 在该代码中，T可以是任何值，包括可能不为nil的值。
		//return nil, false

		var zero T
		return zero, false
	}
}

//Zero 在该代码中，T可以是任何值，包括可能不为nil的值。
// 我们可以利用var语句来解决这个问题，它生成一个新变量，并将其初始化为该类型的零值：
func Zero[T any](T) T {
	var zero T
	return zero
}
