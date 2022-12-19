package pkg

import (
	"github.com/zqlpaopao/tool/struct-tool/pkg"
	"strconv"
	"testing"
)

func BenchmarkPopulateStruct(b *testing.B){

	for i := 0; i < b.N; i++ {
		var m People

		if err := populateStruct(&m);err != nil{
			b.Fatal(err)
		}
		if m.Name != "zhangSan"{
			b.Fatal("name is err")
		}
	}
}
func BenchmarkPopulateStructCache(b *testing.B){

	for i := 0; i < b.N; i++ {
		var m People

		if err := populateStructCache(&m);err != nil{
			b.Fatal(err)
		}
		if m.Name != "zhangSan"{
			b.Fatal("name is err")
		}
	}
}
func BenchmarkPopulateStructUnsafe(b *testing.B){

	for i := 0; i < b.N; i++ {
		var m People

		if err := populateStructUnsafe(&m);err != nil{
			b.Fatal(err)
		}
		if m.Name != "zhangSan"{
			b.Fatal("name is err")
		}
	}
}
func BenchmarkPopulateStructUintPtr(b *testing.B){

	for i := 0; i < b.N; i++ {
		var m People

		if err := populateStructUintPtr(&m);err != nil{
			b.Fatal(err)
		}
		if m.Name != "zhangSan"{
			b.Fatal("name is err")
		}
	}
}


func BenchmarkSimplePopulateStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var m People

		simplePopulateStruct(&m)
		if m.Name != "zhangSan"{
			b.Fatal("name is err")
		}
	}
}

type Model struct {
	Str string `json:"str"`
	//StrPtr *string `json:"str_ptr"`
	Int8   int8 `json:"int_8"`
	//Int8Ptr *int8 `json:"int_8_ptr"`
	Slice []string `json:"slice"`
	//SlicePtr *[]string `json:"slice_ptr"`
	Float32 float32 `json:"float_32"`
	//Float32Ptr *float32 `json:"float_32_ptr"`
}

func BenchmarkMap2Struct(b *testing.B) {

	int8Ptr := 18
	slice :=  []string{
		"slice1","slice2",
	}
	slicePtr := []string{
		"slice1Ptr","slice2Ptr",
	}
	strPtr := "str_ptrs"

	float32Ptr := 32.01
	var m = map[string]interface{}{
		"str" :"str",
		"str_ptr" :&strPtr,
		"int_8" :int8(8),
		"int_8_ptr" :&int8Ptr,
		"slice" :slice,
		"slice_ptr" :&slicePtr,
		"float_32" :32.00,
		"float_32_ptr" :&float32Ptr,
	}
	for i := 0; i < b.N; i++ {
		m["str"] = strconv.Itoa(i)
		var structs Model
		if err := Map2Struct(m,&structs,"json");nil != err{
			b.Fatal(err)
		}
		//fmt.Println(structs)
		//fmt.Println(structs.StrPtr)
		//fmt.Println(structs.StrPt)
	}
}
func BenchmarkMap2StructOld(b *testing.B) {

	int8Ptr := 18
	slice :=  []string{
		"slice1","slice2",
	}
	slicePtr := []string{
		"slice1Ptr","slice2Ptr",
	}
	strPtr := "str_ptrs"

	float32Ptr := 32.01
	var m = map[string]interface{}{
		"str" :"str",
		"str_ptr" :&strPtr,
		"int_8" :int8(8),
		"int_8_ptr" :&int8Ptr,
		"slice" :slice,
		"slice_ptr" :&slicePtr,
		"float_32" :32.00,
		"float_32_ptr" :&float32Ptr,
	}
	for i := 0; i < b.N; i++ {
		m["str"] = strconv.Itoa(i) +"----"
		var structs Model
		//pkg.MapToStruct()
		if err := pkg.MapToStruct(m,&structs,"json");nil != err{
			b.Fatal(err)
		}
		//fmt.Println(structs.Str)
		//fmt.Println(structs.StrPtr)
		//fmt.Println(structs.StrPt)
	}
}


func BenchmarkMap2StructOverPopulateStructDescriptor(b *testing.B){
	int8Ptr := 18
	slice :=  []string{
		"slice1","slice2",
	}
	slicePtr := []string{
		"slice1Ptr","slice2Ptr",
	}
	strPtr := "str_ptrs"

	float32Ptr := 32.01
	var m = map[string]interface{}{
		"str" :"str",
		"str_ptr" :&strPtr,
		"int_8" :int8(8),
		"int_8_ptr" :&int8Ptr,
		"slice" :slice,
		"slice_ptr" :&slicePtr,
		"float_32" :32.00,
		"float_32_ptr" :&float32Ptr,
	}
	m = m
	descriptor, err := describeType((*Model)(nil))
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		var m Model
		if err := populateStructUintDescriptor(&m,descriptor);err != nil{
			b.Fatal(err)
		}
		if m.Str != "zhangSan"{
			b.Fatal("name is err")
		}
	}
}


func BenchmarkMap2StructOver(b *testing.B) {
	int8Ptr := 18
	slice :=  []string{
		"slice1","slice2",
	}
	slicePtr := []string{
		"slice1Ptr","slice2Ptr",
	}
	strPtr := "str_ptrs"

	float32Ptr := 32.01

	var m = map[string]interface{}{
		"str" :"str",
		"str_ptr" :&strPtr,
		"int_8" :int8(8),
		"int_8_ptr" :&int8Ptr,
		"slice" :slice,
		"slice_ptr" :&slicePtr,
		"float_32" :float32(32.00),
		"float_32_ptr" :&float32Ptr,
	}
	//
	//m = m

	var tagNames = []*TagName{
		{
			StructName: "Str",
			Type: StringType,
			MapKey:  "str",
		},
		{
			StructName: "Int8",
			Type: Int8Type,
			MapKey:  "int_8",
		},
		{
			StructName: "Slice",
			Type: SliceStringType,
			MapKey:  "slice",
		},
		{
			StructName: "Float32",
			Type: Float32Type,
			MapKey:  "float_32",
		},
		{
			StructName: "Slice",
			Type: SliceStringType,
			MapKey:  "slice",
		},

	}

	b.StartTimer()
	err := DescribeStructUnsafePointer((*Model)(nil),tagNames)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		m["str"] = strconv.Itoa(i+100)
		var structs Model
		//pkg.MapToStruct()
		if err := Map2StructOver(&structs,tagNames,m);nil != err{
			b.Fatal(err)
		}
		//fmt.Println(structs)
		//fmt.Println(structs.StrPtr)
		//fmt.Println(structs.StrPt)
	}
}
