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
