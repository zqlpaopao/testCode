package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name *string `json:"name,omitempty"`
	Age *int64 `json:"age,omitempty"`
	Sex *int64 `json:"sex,omitempty"`
}

func main(){
	var a = `{"age":15,"sex":1,"name":""}`
	req := Person{}
	_ = json.Unmarshal([]byte(a),&req)

	fmt.Printf("%+v\n",*req.Name)
	fmt.Printf("%d\n",*req.Age)
}
