# 1、安装
```go
go get github.com/go-playground/validator/v10
```

# 2、使用样例
```go
/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/4/9 下午1:24
 */


package main

import (
"fmt"
"github.com/go-playground/validator/v10"
)

// 用户信息
type User struct {
	Name      string     `validate:"required"`
	Age       uint8      `validate:"gte=0,lte=100"`
	Email     string     `validate:"required,email"`
	Phone     int64      `validate:"required,number"`
	Addresses []*Address `validate:"required,dive,required"`
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
	//创建一个示例
	validate = validator.New()

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
	validateStruct(user)
	//验证某一单一变量
	validateVariable()
}

//-- ----------------------------
//--> @Description 结构体验证
//--> @Param
//--> @return
//-- ----------------------------
func validateStruct(user *User) {
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
			fmt.Println()
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

```

# 3、自定义函数验证

```go
err := validate.RegisterValidation("checkName", checkName)
```

```go
package main

import (
    "fmt"
    "gopkg.in/go-playground/validator.v9"
    "unicode/utf8"
)

// User contains user information
type UserInfo struct {
    Name           string     `validate:"checkName"`
    Number float32 `validate:"numeric"`
}
// 自定义验证函数
func checkName(fl validator.FieldLevel) bool {
    count := utf8.RuneCountInString(fl.Field().String())
    fmt.Printf("length: %v \n", count)
    if  count > 5 {
        return false
    }
    return true
}

func main() {
    validate := validator.New()
        //注册自定义函数，与struct tag关联起来
    err := validate.RegisterValidation("checkName", checkName)
    user := &UserInfo{
        Name:            "我是中国人，我爱自己的祖国",
        Number:         23,
    }
    err = validate.Struct(user)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            fmt.Println(err)
        }
        return
    }
    fmt.Println("success")
}
```

# 4、返回错误为中文

重要的点

中文翻译插件一定和对应的validator版本对应

```go
//中文翻译器
	zh_ch := zh.New()
	uni := ut.New(zh_ch)
	trans, _ := uni.GetTranslator("zh")

	//创建一个示例
	validate = validator.New()
	//验证器注册翻译器
	err := zhTranslations.RegisterDefaultTranslations(validate, trans)

fmt.Println("以下是中文错误翻译")
	//中文错误翻译
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			//翻译错误信息
			fmt.Println(err.Translate(trans))
		}
		return
	}
```



```go
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

```



# 5、对应tag配置

https://pkg.go.dev/gopkg.in/go-playground/validator.v10

对应的使用看官网的使用介绍

###  比较运算符

| 运算符 | 运算描述 |
| :----- | :------- |
| eq     | 等于     |
| gt     | 大于     |
| gte    | 大于等于 |
| lt     | 小于     |
| lte    | 小于等于 |
| ne     | 不等于   |

### 1. 字段验证运算

| 运算符        | 运算描述                 |
| :------------ | :----------------------- |
| eqcsfield     | 跨不同结构体字段相等     |
| eqfield       | 同一结构体字段相等       |
| fieldcontains | 包含字段                 |
| fieldexcludes | 未包含字段               |
| gtcsfield     | 跨不同结构体字段大于     |
| gtecsfield    | 跨不同结构体字段大于等于 |
| gtefield      | 同一结构体字段大于等于   |
| gtfield       | 同一结构体字段相等       |
| ltcsfield     | 跨不同结构体字段小于     |
| ltecsfield    | 跨不同结构体字段小于等于 |
| ltefield      | 同一个结构体字段小于等于 |
| ltfield       | 同一个结构体字段小于     |
| necsfield     | 跨不同结构体字段不想等   |
| nefield       | 同一个结构体字段不想等   |

### 2. 网络字段验证运算

| 运算符        | 运算描述                        |
| :------------ | :------------------------------ |
| cidr          | 有效 CIDR                       |
| cidrv4        | 有效 CIDRv4                     |
| cidrv6        | 有效 CIDRv6                     |
| datauri       | 是否有效 URL                    |
| fqdn          | 有效完全合格的域名 (FQDN)       |
| hostname      | 是否站点名称 RFC 952            |
| hostname_port | 是否站点端口                    |
| ip            | 是否包含有效 IP                 |
| ip4_addr      | 是否有效 IPv4                   |
| ip6_addr      | 是否有效 IPv6                   |
| ip_addr       | 是否有效 IP                     |
| ipv4          | 是否有效 IPv4                   |
| ipv6          | 是否有效 IPv6                   |
| mac           | 是否媒体有效控制有效 MAC 地址   |
| tcp4_addr     | 是否有效 TCPv4 传输控制协议地址 |
| tcp6_addr     | 是否有效 TCPv6 传输控制协议地址 |
| tcp_addr      | 是否有效 TCP 传输控制协议地址   |
| udp4_addr     | 是否有效 UDPv4 用户数据报协议地 |
| udp6_addr     | 是否有效 UDPv6 用户数据报协议地 |
| udp_addr      | 是否有效 UDPv 用户数据报协议地  |
| unix_addr     | Unix域套接字端点地址            |
| uri           | 是否包含有效的 URI              |
| url           | 是否包含有效的 URL              |

## 3. 字符串验证运算

| 运算符          | 运算描述                        |
| :-------------- | :------------------------------ |
| alpha           | 是否全部由字母组成的            |
| alphanum        | 是否全部由数字组成的            |
| alphanumunicode | 是否全部由 unicode 字母数字组成 |
| alphaunicode    | 是否全部由 unicode 字母组成     |
| ascii           | ASCII                           |
| contains        | 是否包含全部                    |
| endswith        | 尾部是否以此字符结束            |
| lowercase       | 是否全小写字符组成              |
| multibyte       | 是否多字节字符                  |
| number          | 是否数字                        |
| numeric         | 是否包含基本的数值              |
| printascii      | 是否可列印的 ASCII              |
| startswith      | 开头是否以此字符开始            |
| uppercase       | 是否全大写字符组成              |

## 4. 数据格式验证运算

| 运算符          | 运算描述                 |
| :-------------- | :----------------------- |
| base64          | 是否 Base64              |
| base64url       | 是否 Base64UR            |
| btc_addr        | 是否 Bitcoin 地址        |
| btc_addr_bech32 | 是否 Bitcoin Bech32 地址 |
| datetime        | 是否有效 Datetime        |
| email           | 是否有效 E-mail          |
| html            | 是否有效 HTML 标签       |
| json            | 是否有效 JSON            |
| rgb             | 是否有效 RGB             |
| rgba            | 是否有效 RGBA            |
| uuid            | 是否有效通用唯一标识符   |

## 5. 其他验证运算

| 运算符    | 运算描述         |
| :-------- | :--------------- |
| dir       | 是否有效目录     |
| file      | 是否有效文件目录 |
| isdefault | 是默认值         |
| len       | 指定长度         |
| max       | 最大值           |
| min       | 最小值           |
| required  | 必须传入         |
| unique    | 唯一的           |