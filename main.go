package main

import (
	"github.com/zqlpaopao/tool/zap-log/src"
)

/*
	提供日志分割和日志保存周期控制
*/

func main(){
	type str struct{
		name string
		age int
		sex []int
	}
	s := str{
		name: "name",
		age:  18,
		sex:  []int{1,2,3,4},
	}
	//debug info 是一个级别 warn和errorshi 是一个级别，不同级别可分别记录
	//debug info 是一个级别 warn和errorshi 是一个级别，不同级别可分别记录
	src.InitLoggerHandler(&src.LogConfig{
		InfoPathFileName: "./demo.log",
		WarnPathFileName: "./demo.log",
		//WithRotationTime: //最大旋转时间 默认值1小时
		//WithMaxAge: //日志最长保存时间，乘以小时 默认禁用
		//WithRotationCount: //保存的最大文件数 //默认禁用
	})
	src.Info("Info",s).Msg("Info")
	src.Warn("Warn",s).Msg("Warn")
	src.Error("Error",s).Msg("Error")
	src.Debug("Debug",s).Msg("Debug")

	src.Warn("Warn",s).Msg("")

}