/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/2/25 上午11:31
 */

package main

import "test/sync/pool/src"

func main(){
	src.Pool()
}

//go:linkname src_runtimec src.runtimec
//go:nosplit
func src_runtimec(){

}