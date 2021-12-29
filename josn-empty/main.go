package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name *string `string-byte:"name,omitempty"`
	Age *int64 `string-byte:"age,omitempty"`
	Sex *int64 `string-byte:"sex,omitempty"`
}

func main(){
	var a = `{"age":15,"sex":1,"name":""}`
	req := Person{}
	_ = json.Unmarshal([]byte(a),&req)

	fmt.Printf("%+v\n",*req.Name)
	fmt.Printf("%d\n",*req.Age)
}
