package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

func main(){
	read()

}

func read(){
	excelFileName := "/Users/zhangsan/Documents/GitHub/testCode/excel/tealeg-xlsx/海南业绩检查.xlsx"
	f, err := xlsx.OpenFile(excelFileName)
	if err != nil {
	fmt.Println(err)
	}


	//fmt.Printf("%#v\n",f.Sheets)
	//fmt.Printf("%#v\n",f.Sheet)
	//fmt.Printf("%#v\n",f.DefinedNames)
	style := xlsx.NewStyle()
	style.Font.Color = xlsx.RGB_White
	style.Fill.BgColor = xlsx.RGB_Dark_Green
	style.Alignment = xlsx.Alignment{
		Horizontal:  "center",
		Vertical:   "center",
	}

	for i1 , v:= range f.Sheets{

		for i ,v1 := range v.Rows{

			if i == 0{
				for i2 , _ := range v1.Cells{
					f.Sheets[i1].Rows[i].Cells[i2].SetStyle(style)

				}
				continue
			}

			f.Sheets[i1].Rows[i].Cells = append(f.Sheets[i1].Rows[i].Cells,&xlsx.Cell{
				Row:            nil,
				Value:          "jkklllllll",
				NumFmt:         "",
				Hidden:         false,
				HMerge:         0,
				VMerge:         0,
				DataValidation: nil,
			})
			//fmt.Println(v1)

		}


	}


	f.Save("/Users/zhangsan/Documents/GitHub/testCode/excel/tealeg-xlsx/海南业绩检查1.xlsx")

		//for _, sheet := range xlFile.Sheets {
		//	for _, row := range sheet.Rows {
		//		for _, cell := range row.Cells {
		//			text := cell.String()
		//			fmt.Printf("%s\n", text)
		//		}
		//	}
		//}

}