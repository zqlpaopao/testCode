package main

import (
	"bufio"
	"fmt"
	"io"
	//"math"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	/*
	文件数据样例
	{"remark": "来电时间：  2021/04/15 13:52:07客户电话：13913xx39xx ", "no": "600020510132021101310210547639", "title": "b-ae0e-0242ac100907", "call_in_date": "2021-04-15 13:52:12", "name": "张三", "_date": "2021-06-15", "name": "张三", "meet": "1"}

	1、我们取出 call_in_date": "2021-04-15 13:52:1的数据写入另一个文件

	*/

	var (
		s time.Time //当前时间
		file *os.File
		fileStat os.FileInfo
		err error
		lastLineSize int64
	)


	s = time.Now()

	if file, err = os.Open("/Users/zhangsan/Downloads/log.txt");err != nil{
		fmt.Println(err)
	}

	defer func() {
		err = file.Close() //close after checking err
	}()

	//queryStartTime, err := time.Parse("2006-01-02T15:04:05.0000Z", startTimeArg)
	//if err != nil {
	//	fmt.Println("Could not able to parse the start time", startTimeArg)
	//	return
	//}
	//
	//queryFinishTime, err := time.Parse("2006-01-02T15:04:05.0000Z", finishTimeArg)
	//if err != nil {
	//	fmt.Println("Could not able to parse the finish time", finishTimeArg)
	//	return
	//}

	/**
	* {name:"log.log", size:911100961, mode:0x1a4,
	modTime:time.Time{wall:0x656c25c, ext:63742660691,
	loc:(*time.Location)(0x1192c80)}, sys:syscall.Stat_t{Dev:16777220,
	Mode:0x81a4, Nlink:0x1, Ino:0x118cba7, Uid:0x1f5, Gid:0x14, Rdev:0,
	Pad_cgo_0:[4]uint8{0x0, 0x0, 0x0, 0x0}, Atimespec:syscall.Timespec{Sec:1607063899, Nsec:977970393},
	Mtimespec:syscall.Timespec{Sec:1607063891, Nsec:106349148}, Ctimespec:syscall.Timespec{Sec:1607063891,
	Nsec:258847043}, Birthtimespec:syscall.Timespec{Sec:1607063883, Nsec:425808150},
	Size:911100961, Blocks:1784104, Blksize:4096, Flags:0x0, Gen:0x0, Lspare:0, Qspare:[2]int64{0, 0}}
	*
	*/
	if fileStat, err = file.Stat();err != nil {
		return
	}

	fileSize := fileStat.Size()//72849354767
	offset := fileSize - 1

	//检测是不是都是空行 只有\n
	for {
		var (
			b []byte
			n int
			char string
		)
		b = make([]byte, 1)
		//从指定位置读取
		if n, err = file.ReadAt(b, offset);err != nil {
			fmt.Println("Error reading file ", err)
			break
		}
		char = string(b[0])
		if char == "\n" {
			break
		}
		offset--
		//获取一行的大小
		lastLineSize += int64(n)
	}

	var (
		lastLine []byte
		logSlice []string
		logSlice1 []string
	)
	//初始化一行大小的空间
	lastLine = make([]byte, lastLineSize)
	_, err = file.ReadAt(lastLine, offset)

	if err != nil {
		fmt.Println("Could not able to read last line with offset", offset, "and lastline size", lastLineSize)
		return
	}
	//根据条件进行区分
	logSlice = strings.Split(strings.Trim(string(lastLine),"\n"),"next_pay_date")
	logSlice1  = strings.Split(logSlice[1],"\"")
	if logSlice1[2] == "2021-06-15"{
		Process(file)

	}

	fmt.Println("\nTime taken - ", time.Since(s))

		fmt.Println(err)
}

func Process(f *os.File) error {

	//读取数据的key，减小gc压力
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 250*1024)
		return lines
	}}

	//读取回来的数据池
	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}

	//一个文件对象本身是实现了io.Reader的 使用bufio.NewReader去初始化一个Reader对象，存在buffer中的，读取一次就会被清空
	r := bufio.NewReader(f) //

	//设置读取缓冲池大小 默认16
	r = bufio.NewReaderSize(r,250 *1024)

	var wg sync.WaitGroup

	for {
		buf := linesPool.Get().([]byte)
		//读取Reader对象中的内容到[]byte类型的buf中
		n, err := r.Read(buf)
		buf = buf[:n]


		if n == 0 {
			if err != nil {
				fmt.Println(err)
				break
			}
			if err == io.EOF {
				break
			}
			return err
		}

		//补齐剩下没满足的剩余
		nextUntillNewline, err := r.ReadBytes('\n')
		//fmt.Println(string(nextUntillNewline))
		if err != io.EOF {
			buf = append(buf, nextUntillNewline...)
		}
		wg.Add(1)
		go func() {
			ProcessChunk(buf, &linesPool, &stringPool)
			wg.Done()
		}()

	}

	wg.Wait()
	return nil
}

func ProcessChunk(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool) {

	//做相应的处理

	////fmt.Println(chunk)
	//
	////fmt.Println(linesPool)
	////fmt.Println(stringPool)
	////var wg2 sync.WaitGroup
	//
	////取出刚才存进去的
	//logs := stringPool.Get().(string)
	//logs = string(chunk)
	//
	////放入申请好的内存中
	////linesPool.Put(chunk)
	//
	////按照行 分切片
	//logsSlice := strings.Split(logs, "\n")
	//
	////放入申请好的地址中
	//stringPool.Put(logs)
	//
	//
	//chunkSize := 300
	//n := len(logsSlice)//切片大小--也就是存在多少行
	//noOfThread := n / chunkSize
	//
	//if n%chunkSize != 0 {
	//	noOfThread++
	//}

	//for i := 0; i < (noOfThread); i++ {
	//
	//	wg2.Add(1)
	//	go func(s int, e int) {
	//		defer wg2.Done() //to avaoid deadlocks
	//		for i := s; i < e; i++ {
	//			fmt.Println(i)
	//			text := logsSlice[i]
	//			if len(text) == 0 {
	//				continue
	//			}
	//			logSlice := strings.SplitN(text, ",", 2)
	//			logCreationTimeString := logSlice[0]
	//
	//			logCreationTime, err := time.Parse("2006-01-02T15:04:05.0000Z", logCreationTimeString)
	//			if err != nil {
	//				fmt.Printf("\n Could not able to parse the time :%s for log : %v", logCreationTimeString, text)
	//				return
	//			}
	//
	//			if logCreationTime.After(start) && logCreationTime.Before(end) {
	//				fmt.Println(text)
	//			}
	//		}
	//
	//	}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice)))))
	//}
	//
	//wg2.Wait()
	//logsSlice = nil
}