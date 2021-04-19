/**
 * @Author: zHangSan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/4/16 下午4:44
 */

package main

import "runtime"

func main(){
	//goGC()

	//gOTraceBack() 查看crash的时候相信信息

	goDebugGc()

}
/*****************************GODEBUG-异常发生打印详细信息****************************************/
func goDebugGc(){
	gc()
}





/*****************************GOTRACEBACK-异常发生打印详细信息****************************************/
func gOTraceBack(){
	panic("kerboom")

}






/*****************************goGC--修改GC频率****************************************/
func goGC(){
	gc()
	runtime.GC()
	gc()
	gc()
	gc()
}

func gc(){
	for i := 0; i< 10000;i++{
		g := make(map[string]string)
		g = g
	}
}

SCHED 0ms: gomaxprocs=4 idleprocs=2 threads=4 spinningthreads=1 idlethreads=0 runqueue=0 gcwaiting=0 nmidlelocked=0 stopwait=0 sysmonwait=0
P0: status=1 schedtick=0 syscalltick=0 m=0 runqsize=0 gfreecnt=0 timerslen=0
P1: status=1 schedtick=0 syscalltick=0 m=3 runqsize=0 gfreecnt=0 timerslen=0
P2: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0 timerslen=0
P3: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0 timerslen=0
M3: p=1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
M2: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
M1: p=-1 curg=17 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=false lockedg=17
M0: p=0 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=1
G1: status=1() m=-1 lockedm=0
G17: status=6() m=1 lockedm=1
G2: status=1() m=-1 lockedm=-1
# command-line-arguments
SCHED 0ms: gomaxprocs=4 idleprocs=1 threads=4 spinningthreads=1 idlethreads=0 runqueue=0 [1 0 0 0]
# command-line-arguments
SCHED 0ms: gomaxprocs=4 idleprocs=1 threads=5 spinningthreads=1 idlethreads=0 runqueue=0 gcwaiting=0 nmidlelocked=1 stopwait=0 sysmonwait=0
P0: status=1 schedtick=0 syscalltick=0 m=3 runqsize=0 gfreecnt=0 timerslen=0
P1: status=1 schedtick=1 syscalltick=0 m=2 runqsize=0 gfreecnt=0 timerslen=0
P2: status=1 schedtick=0 syscalltick=0 m=4 runqsize=0 gfreecnt=0 timerslen=0
P3: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0 timerslen=0
M4: p=2 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
M3: p=0 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
M2: p=1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=true blocked=false lockedg=-1
M1: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
M0: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=true lockedg=1
G1: status=1(chan receive) m=-1 lockedm=0
G2: status=4(force gc (idle)) m=-1 lockedm=-1
G3: status=1() m=-1 lockedm=-1
G4: status=4(GC scavenge wait) m=-1 lockedm=-1
SCHED 0ms: gomaxprocs=4 idleprocs=1 threads=5 spinningthreads=1 idlethreads=0 runqueue=0 gcwaiting=0 nmidlelocked=1 stopwait=0 sysmonwait=0
P0: status=1 schedtick=0 syscalltick=0 m=3 runqsize=0 gfreecnt=0 timerslen=0
P1: status=1 schedtick=2 syscalltick=0 m=2 runqsize=0 gfreecnt=0 timerslen=0
P2: status=1 schedtick=0 syscalltick=0 m=4 runqsize=0 gfreecnt=0 timerslen=0
P3: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0 timerslen=0
M4: p=2 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=true blocked=false lockedg=-1
M3: p=0 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=true blocked=false lockedg=-1
M2: p=1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
M1: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 spinning=false blocked=false lockedg=-1
M0: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=true lockedg=1
G1: status=1(chan receive) m=-1 lockedm=0
G2: status=4(force gc (idle)) m=-1 lockedm=-1
G3: status=4(GC sweep wait) m=-1 lockedm=-1
G4: status=4(GC scavenge wait) m=-1 lockedm=-1