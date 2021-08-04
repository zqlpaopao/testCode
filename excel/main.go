package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

func StreamWriterFunc(contents [][]string)  {
	//打开工作簿
	file, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		return
	}
	sheet_name := "Sheet1"
	//获取流式写入器
	streamWriter, _ := file.NewStreamWriter(sheet_name)
	if err != nil {
		fmt.Println(err)
	}

	rows, _ := file.GetRows(sheet_name)	//获取行内容
	cols, _ := file.GetCols(sheet_name)	//获取列内容
	fmt.Println("行数rows:  ", len(rows),"列数cols:  ", len(cols))

	//将源文件内容先写入excel
	for rowid , row_pre:= range rows{
		row_p := make([]interface{}, len(cols))
		for colID_p := 0; colID_p < len(cols); colID_p++ {
			//fmt.Println(row_pre)
			//fmt.Println(colID_p)
			if row_pre == nil {
				row_p[colID_p] = nil
			}else {
				row_p[colID_p] = row_pre[colID_p]
			}
		}
		cell_pre, _ := excelize.CoordinatesToCellName(1, rowid+1)
		if err := streamWriter.SetRow(cell_pre, row_p); err != nil {
			fmt.Println(err)
		}
	}

	//将新加contents写进流式写入器
	for rowID := 0; rowID < len(contents); rowID++ {
		row := make([]interface{}, len(contents[0]))
		for colID := 0; colID < len(contents[0]); colID++ {
			row[colID] = contents[rowID][colID]
		}
		cell, _ := excelize.CoordinatesToCellName(1, rowID+len(rows)+1) //决定写入的位置
		if err := streamWriter.SetRow(cell, row); err != nil {
			fmt.Println(err)
		}
	}

	//结束流式写入过程
	if err := streamWriter.Flush(); err != nil {
		fmt.Println(err)
	}
	//保存工作簿
	if err := file.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

