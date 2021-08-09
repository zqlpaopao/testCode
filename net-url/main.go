package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/url"
	"reflect"
)




/************************************** 字符床特殊字符 @、= 、！相关******************/
func SpecialChar(){
	sessid := ""
	var err error
	if sessid,err = sessionId();nil != err{
		panic(err)
	}
	fmt.Println(sessid)
	//uPvjVSMl6eksU9acFeV_iFrFBgOTvR1D-6u69PTtz2A=


	//将字符串中的特殊字符转换 = -----> %3D
	encodeSessId := url.QueryEscape(sessid)
	fmt.Println(encodeSessId)
	//uPvjVSMl6eksU9acFeV_iFrFBgOTvR1D-6u69PTtz2A%3D

	//将字符传中转换的特殊字符还原
	deCodeSessId ,errs := url.QueryUnescape(encodeSessId)
	fmt.Println(errs)
	fmt.Println(deCodeSessId)
	//uPvjVSMl6eksU9acFeV_iFrFBgOTvR1D-6u69PTtz2A=
}

/**
	对字符串的特殊字符进行转换和反转换
	例
	= ------》 %3D
	%3D -----》 =
**/
func sessionId()(str string,err error){
	b := make([]byte,32)

	//ReadFull 从 rand.Reader 精确地读取 len(b) 字节数据填充进 b
	//rand.Reader 是一个全局、共享的密码用强随机数生成器
	if _,err = io.ReadFull(rand.Reader,b);nil != err{
		return "",err
	}

	fmt.Println(b)//[71 148 194 155 86 128 235 223 185 2 177 179 52 164 167 244 100 8 210 64 135 70 87 228 65 17 16 181 80 136 203 103]

	return base64.URLEncoding.EncodeToString(b),err
}

/************************************** URL格式相关******************/
/**
	type URL struct {
    Scheme   string    //具体指访问服务器上的资源使用的哪种协议
    Opaque   string    // 编码后的不透明数据
    User     *Userinfo // 用户名和密码信息，有些协议需要传入明文用户名和密码来获取资源，比如 FTP
    Host     string    // host 或 host:port，服务器地址，可以是 IP 地址，也可以是域名信息
    Path     string  //路径，使用"/"分隔
    RawQuery string // 编码后的查询字符串，没有'?'
    Fragment string // 引用的片段（文档位置），没有'#'
}
*/
func URLFormat(){
	u := url.URL{
		Scheme:      "",
		Opaque:      "",
		User:        nil,
		Host:        "example.com",
		Path:        "foo",
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	fmt.Println(u.IsAbs())//false 检测是否有协议
	u.Scheme ="uil"
	fmt.Println(u.IsAbs())//true 只要不为空

	us := &url.URL{
		Scheme:   "https",
		User:     url.UserPassword("me", "pass"),
		Host:     "example.com",
		Path:     "foo/bar",
		RawQuery: "x=1&y=2",
		Fragment: "anchor",
	}

	//对参数进行解析
	fmt.Println(us.Query()) //map[x:[1] y:[2]]

	//解析get请求参数
	uc, err := url.Parse("https://example.org/path?foo=bar")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(uc.RequestURI()) //    /path?foo=bar
	fmt.Println(u.RequestURI()) // 		foo
	fmt.Println(us.RequestURI()) //     foo/bar?x=1&y=2
	fmt.Println(us.RawQuery) //     	x=1&y=2
}

/************************************** 将URL构造合法url******************/
func makeUrl(){
	u := &url.URL{
		Scheme:   "https",
		User:     url.UserPassword("me", "pass"),
		Host:     "example.com",
		Path:     "foo/bar",
		RawQuery: "x=1&y=2",
		Fragment: "anchor",
	}
	//这是第一种形式
	//fmt.Println(u.String()) //https://me:pass@example.com/foo/bar?x=1&y=2#anchor
	u.Opaque = "opaque"//将域名变为不透明
	//这是第二种形式
	//fmt.Println(u.String()) //https:opaque?x=1&y=2#anchor



	//将特殊字符转移
	up ,err := url.Parse("https://example.com/foo/bar  anchor")
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(up.EscapedPath())///foo/bar%20%20anchor
}

/**************************************解析url参数******************/
func parseUrlArgs(){
	u, err := url.Parse("https://example.org:8000/path")//IPV4
	if err != nil {
		log.Fatal(err)
	}
	//获取域名
	fmt.Println(u.Hostname()) //example.org

	//获取端口
	fmt.Println(u.Port())
	fmt.Println(u.Path)
	fmt.Println(u.String())

	u, err = url.Parse("https://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:17000") //IPV6
	if err != nil {
		log.Fatal(err)
	}
	//获取域名
	fmt.Println(u.Hostname())//2001:0db8:85a3:0000:0000:8a2e:0370:7334

	//获取端口
	fmt.Println(u.Port())
}

/**************************************生成ref******************/
func ref(){
	base, err := url.Parse("http://example.com/directory/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(base)//http://example.com/directory/
	result, err := base.Parse("./search?q=dotnet")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)//http://example.com/directory/search?q=dotnet

	//返回相对ref
	u, err := url.Parse("../../..//search?q=dotnet")//相对路径的不同会影响返回的结果
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u)//     ../../..//search?q=dotnet

	base, err = url.Parse("http://example.com/directory/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(base)//http://example.com/directory/
	fmt.Println(base.ResolveReference(u))//http://example.com/search?q=dotnet
}

/**************************************二进制转换******************/
func binary(){
	u, _ := url.Parse("https://example.org")
	b, err := u.MarshalBinary() //将其转成二进制
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reflect.TypeOf(b)) //[]uint8
	fmt.Println(b) //[104 116 116 112 115 58 47 47 101 120 97 109 112 108 101 46 111 114 103]
	fmt.Printf("%s\n", b) //https://example.org



	us := &url.URL{}
	//将其从二进制转成 url.URL 类型
	err = u.UnmarshalBinary([]byte("https://example.org:8000/foo"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reflect.TypeOf(us)) //*url.URL
	fmt.Println(u) //https://example.org:8000/foo
	fmt.Println(u.Hostname()) //example.org
	fmt.Println(u.Port()) //8000
}

/**************************************userInfo******************/
func userInfo(){
	u := &url.URL{
		Scheme:   "https",
		User:     url.UserPassword("me", "pass"),
		Host:     "example.com",
		Path:     "foo/bar",
		RawQuery: "x=1&y=2",
		Fragment: "anchor",
	}
	fmt.Println(u.User.Username()) //me
	password, b := u.User.Password()
	if b == false{
		log.Fatal("can not get password")
	}
	fmt.Println(password) //pass
	fmt.Println(u.User.String()) //me:pass

}
/**************************************获取url信息转换为map******************/
func urlArgs(){
	v, err := url.ParseQuery("friend=Jess&friend=Sarah&fruit=apple&name=Ava")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v) //map[friend:[Jess Sarah] fruit:[apple] name:[Ava]]
}

/**************************************增加、删除、获取url的参数值信息******************/
func getSetDelQuery(){
	v := url.Values{}
	v.Set("name", "Ava")
	v.Add("friend", "Jess")
	v.Add("friend", "Sarah")
	v.Add("fruit", "apple")

	fmt.Println(v.Get("name"))
	fmt.Println(v.Get("friend"))
	fmt.Println(v["friend"])
	fmt.Println(v.Encode())

	v.Del("name")
	fmt.Println(v.Encode())
}

func main(){
	//字符串获取url参数中的特殊字符转换
	//SpecialChar()

	//URL及参数解析各项相关
	//URLFormat()

	//构造请求url
	//makeUrl()

	//解析url参数
	//parseUrlArgs()

	//生成ref，相对路径
	//ref()

	//二进制操作
	//binary()

	//解析用户信息
	//userInfo()

	//增加、删除、获取url键值对信息
	getSetDelQuery()
}