///**
// * @Author: zhangsan
// * @Description:
// * @File:  main
// * @Version: 1.0.0
// * @Date: 2021/4/9 下午1:24
// */
//
//
package main

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/locales/zh"
	//enTranslations "github.com/go-playground/validator/v10/translations/en"//英文
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"

)

// 用户信息
type User struct {
	Name      string     `validate:"required"`
	Age       uint8      `validate:"gte=0,lte=100"`
	Email     string     `validate:"required,email"`
	Phone     int64      `validate:"required,number"`
	Addresses []*Address `validate:"required,dive,required"`//dive-进入切片
}

// 用户地址信息
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Phone  string `validate:"required"`
}

// 用 Validate 的单个实例来缓存结构体信息
var validate *validator.Validate


func main() {
	//中文翻译器
	zh_ch := zh.New()
	uni := ut.New(zh_ch)
	trans, _ := uni.GetTranslator("zh")

	//创建一个示例
	validate = validator.New()
	//验证器注册翻译器
	err := zhTranslations.RegisterDefaultTranslations(validate, trans)
	fmt.Println(err)
	address := &Address{
		Street: "chengDe weiChang.32",
		City:   "chengDe",
		Phone:  "101-101010101",
	}

	user := &User{
		Name:      "zhangsan",
		Age:       180,
		Email:     "zhangsan@gmail.com",
		Phone:     13889898899,
		Addresses: []*Address{address},
	}
	//验证结构体内容
	validateStruct(user,trans)
	//验证某一单一变量
	validateVariable()
}

//-- ----------------------------
//--> @Description 结构体验证
//--> @Param
//--> @return
//-- ----------------------------
func validateStruct(user *User,trans ut.Translator) {
	err := validate.Struct(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)  //Key: 'User.Age' Error:Field validation for 'Age' failed on the 'lte' tag
			fmt.Println(err.Namespace())//User.Age
			fmt.Println(err.Field())//Age
			fmt.Println(err.StructNamespace())//User.Age
			fmt.Println(err.StructField())//Age
			fmt.Println(err.Tag())//lte
			fmt.Println(err.ActualTag())//lte
			fmt.Println(err.Kind())//uint8
			fmt.Println(err.Type())//uint8
			fmt.Println(err.Value())//180
			fmt.Println(err.Param())//180
		}
		//return
	}

	fmt.Println("以下是中文错误翻译")
	//中文错误翻译
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			//翻译错误信息
			fmt.Println(err.Translate(trans))
		}
		return
	}

}

//单一结构体的数据验证
func validateVariable() {
	myEmail := "zhangsan@gmail.com"

	errs := validate.Var(myEmail, "required,email")

	if errs != nil {
		fmt.Println(errs)
		return
	}
}
