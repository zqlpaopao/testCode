/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/1/29 下午6:42
 */

package main

import "fmt"

func main() {
	NodesJob := map[string]string{
		"12345":"work",
		"123455":"work1",
		"1234556":"work1",
	}
	jobSortCopy := []string{"12","34"}
	var avg int64 = 1
	workerJobMap := make(map[string][]string)

	for i := range NodesJob {
		if avg >= int64(len(jobSortCopy)) {
			workerJobMap[i] = jobSortCopy
			break
		} else {
			workerJobMap[i] = jobSortCopy[:avg]
			jobSortCopy = jobSortCopy[avg:]
		}
	}

	fmt.Println(workerJobMap)
}