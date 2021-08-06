package main

import "github.com/tealeg/xlsx"

func main() {
	file, err := xlsx.OpenFile("demo.xlsx")
	if err != nil {
		panic(err.Error())
	}
	file.Sheets[0].Rows[1].Cells[0].Value = "李四"
	err = file.Save("demo.xlsx")
	if err != nil {
		panic(err.Error())
	}
}