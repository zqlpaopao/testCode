package pkg

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)
var cache = make(map[reflect.Type][]int)

type People struct {
	Name string
	Age int
	sex int
}



func simplePopulateStruct(people *People){
	people.Name = "zhangSan"
}

//反射基本版
func populateStruct(people interface{})(err error){
	val := reflect.ValueOf(people)

	if val.Type().Kind() != reflect.Ptr{
		return errors.New("pass in a pointer")
	}
	elm := val.Elem()

	if elm.Type().Kind() != reflect.Struct{
		return errors.New("type is not struct")
	}

	fval := elm.FieldByName("Name")
	fval.SetString("zhangSan")
	return
}

//反射-cache基本版
func populateStructCache(in interface{})(err error){
	typ := reflect.TypeOf(in)
	index , ok := cache[typ]
	if !ok{
		if typ.Kind() != reflect.Ptr{
			return errors.New("pass in a pointer")
		}

		if typ.Elem().Kind() != reflect.Struct{
			return errors.New("type is not struct")
		}

		f, ok1 := typ.Elem().FieldByName("Name")
		if !ok1{
			return errors.New("struct does not have field Name")
		}
		index = f.Index
		cache[typ] = index

	}

	val := reflect.ValueOf(in)
	elm := val.Elem()
	fVal := elm.FieldByIndex(index)
	fVal.SetString("zhangSan")

	return
}

var unsafeCache = make(map[reflect.Type]uintptr)

type modelFace struct {
	typ   unsafe.Pointer
	value unsafe.Pointer
}

//反射-unsafe
func populateStructUnsafe(in interface{})(err error){
	typ := reflect.TypeOf(in)
	offset, ok := unsafeCache[typ]
	if !ok {
		if typ.Kind() != reflect.Ptr {
			return fmt.Errorf("you must pass in a pointer")
		}
		if typ.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("you must pass in a pointer to a struct")
		}
		f, ok := typ.Elem().FieldByName("Name")
		if !ok {
			return fmt.Errorf("struct does not have field B")
		}
		if f.Type.Kind() != reflect.String {
			return fmt.Errorf("field Name should be an string")
		}
		offset = f.Offset
		unsafeCache[typ] = offset
	}
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*string)(unsafe.Pointer(uintptr(structPtr) + offset)) = "zhangSan"
	return
}


var uintPtrCache = make(map[uintptr]uintptr)


//反射-uintptr
func populateStructUintPtr(in interface{})(err error){
	inf := (*modelFace)(unsafe.Pointer(&in))
	offset, ok := uintPtrCache[uintptr(inf.typ)]
	if !ok {
		typ := reflect.TypeOf(in)
		if typ.Kind() != reflect.Ptr {
			return fmt.Errorf("you must pass in a pointer")
		}
		if typ.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("you must pass in a pointer to a struct")
		}
		f, ok := typ.Elem().FieldByName("Name")
		if !ok {
			return fmt.Errorf("struct does not have field B")
		}
		if f.Type.Kind() != reflect.String {
			return fmt.Errorf("field B should be an String")
		}
		offset = f.Offset
		uintPtrCache[uintptr(inf.typ)] = offset
	}

	*(*string)(unsafe.Pointer(uintptr(inf.value) + offset)) = "zhangSan"


	return
}

type typeDescriptor uintptr
//描述符
func populateStructUintDescriptor(in interface{}, ti typeDescriptor)(err error){
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*string)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = "zhangSan"
	return nil
}

func describeType(in interface{}) (typeDescriptor, error) {
	typ := reflect.TypeOf(in)
	if typ.Kind() != reflect.Ptr {
		return 0, fmt.Errorf("you must pass in a pointer")
	}
	if typ.Elem().Kind() != reflect.Struct {
		return 0, fmt.Errorf("you must pass in a pointer to a struct")
	}
	f, ok := typ.Elem().FieldByName("Name")
	if !ok {
		return 0, fmt.Errorf("struct does not have field B")
	}
	if f.Type.Kind() != reflect.String {
		return 0, fmt.Errorf("field B should be an int")
	}
	return typeDescriptor(f.Offset), nil
}

//***************************************修改全部的***************************************//

//***************************************修改全部的***************************************//

type uintDescriptor uintptr

type modelFaces struct {
	typ   unsafe.Pointer
	value unsafe.Pointer
}

func Map2Struct(m,s interface{},tagName string)(err error){

	//如果传递的是指针，直接解引用
	mVal := reflect.Indirect(reflect.ValueOf(m))
	mValType := mVal.Type()
	if mValType.Kind() != reflect.Map {
		return errors.New("you must pass in a map")
	}

	sVal := reflect.Indirect(reflect.ValueOf(s))
	sValType := sVal.Type()
	if sValType.Kind() != reflect.Struct {
		return errors.New("you must pass in a pointer to a struct")
	}


	for i := 0; i < sVal.NumField(); i++ {
		//需要先判断是否有tag标签
		sField := sValType.Field(i)
		var name string
		if tagName != "" {
			name = sField.Tag.Get(tagName)
		}
		if name == "" {
			name = sField.Name
		}
		mKey := mVal.MapIndex(reflect.ValueOf(name))
		if !mKey.IsValid() {
			continue
		}
		if mKey.IsZero() {
			continue
		}
		//由于从map中获取的int值都是默认int，float默认float64，因此需要做特殊处理
		values := reflect.Indirect(mKey.Elem()).Interface()
		//fmt.Println(name,"name")
		//fmt.Println(mKey,"mk")
		//fmt.Println(values, reflect.TypeOf(values))
		switch reflect.TypeOf(values).Kind() {
		case reflect.Int8:
			isIntType(s,sField.Type.Kind(),uintDescriptor(sField.Offset), int64(values.(int8)))

		case reflect.Int:
			isIntType(s,sField.Type.Kind(),uintDescriptor(sField.Offset), int64(values.(int)))
		case reflect.Int32:
			 isIntType(s,sField.Type.Kind(), uintDescriptor(sField.Offset),int64(values.(int32)))

		case reflect.Int64:
			isIntType(s,sField.Type.Kind(), uintDescriptor(sField.Offset),values)

		case reflect.Float64:
			isFloatType(s,sField.Type.Kind(), uintDescriptor(sField.Offset),values)

		case reflect.String :
			isStringType(s,sField.Type.Kind(), uintDescriptor(sField.Offset),values.(string))
		case reflect.Slice:
			isStringType(s,sField.Type.Kind(), uintDescriptor(sField.Offset),values.([]string))

		default:
			if mKey.Elem().Type() != sField.Type {
				continue
			}
		}
	}
	return
}


//判断是否是int大类
func isIntType(s interface{},k reflect.Kind, offset uintDescriptor, values interface{})  {
	switch k {
	case reflect.Uint8:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*uint8)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = uint8(values.(int64))
	case reflect.Uint16:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*uint16)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = uint16(values.(int64))

	case reflect.Uint32:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*uint32)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = uint32(values.(int64))

	case reflect.Uint64:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*uint64)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = uint64(values.(int64))

	case reflect.Int8:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*int8)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = int8(values.(int64))

	case reflect.Int16:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*int16)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = int16(values.(int64))

	case reflect.Int32:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*int32)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = int32(values.(int64))

	case reflect.Int64:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*int64)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = values.(int64)
	case reflect.Float32:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*float32)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = float32(values.(int64))
	case reflect.Float64:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*float64)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = float64(values.(int64))

	default:
	}
}

//判断是否float大类
func isFloatType(s interface{},k reflect.Kind, offset uintDescriptor,values interface{}){
	switch k {
	case reflect.Float32:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*float32)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = float32(values.(float64))
	case reflect.Float64:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*float64)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = values.(float64)

	default:

	}
}


func isStringType(s interface{},k reflect.Kind, offset uintDescriptor,values interface{}){
	switch k {
	case reflect.String:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*string)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = values.(string)
	case reflect.Slice:
		structPtr := (*modelFace)(unsafe.Pointer(&s)).value
		*(*[]string)(unsafe.Pointer(uintptr(structPtr) + uintptr(offset))) = values.([]string)
	default:

	}
}


