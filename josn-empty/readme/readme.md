# 1、json 去除空值
定义结构体的时候加上tag omitempty 标识
```go
package main

import (
	"encoding/json"
	"fmt"
)

type person struct {
	Name string `json:"name,omitempty"`
	Age int64 `json:"age,omitempty"`
	Sex int64 `json:"sex,omitempty"`
}

func main(){
	var (
		per person
		b []byte
		err error
	)

	per = person{
		Name: "",
		Age:  15,
		Sex:  1,
	}

	b , err = json.Marshal(per)

	fmt.Println(err)
	fmt.Println(string(b))
}
```
结果
```go
<nil>
{"age":15,"sex":1}
```
结果中祛去除了name，因为它为空

# 2、就想用空显示
这个时候我不想去除空，可能前端传过来的就是""值
在现有的一些Golang库中，go-github中有体现

核心就是转为指针类型
```go
package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_JSON_Empty(t *testing.T) {
	jsonData := `{"age":12,"name":"xyz","sex":0}`
	req := Person{}
	_ = json.Unmarshal([]byte(jsonData), &req)
	fmt.Printf("%+v\n", req)
	fmt.Printf("%s\n", *req.Name)
	fmt.Printf("%d\n", *req.Sex)
}
func Test_JSON_Nil(t *testing.T) {
	jsonData := `{"id":"1234","name":"xyz","age":0}`
	req := Person{}
	_ = json.Unmarshal([]byte(jsonData), &req)
	fmt.Printf("%+v\n", req)
	fmt.Printf("%d\n", *req.Age)
}

```

结果
```go
go test -v
=== RUN   Test_JSON_Empty
{Name:0xc00008e580 Age:0xc0000a6178 Sex:0xc0000a6188}
xyz
0
--- PASS: Test_JSON_Empty (0.00s)
=== RUN   Test_JSON_Nil
{Name:0xc00008e600 Age:0xc0000a61d0 Sex:<nil>}
0
--- PASS: Test_JSON_Nil (0.00s)
PASS
ok      test/josn-empty 0.005s

```

可以看到没有的sex 为nil 而为0的age 实际是0，不为nil
