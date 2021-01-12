/**
 * @Author: zhangsan
 * @Description:
 * @File:  split_test.go
 * @Version: 1.0.0
 * @Date: 2021/1/11 上午11:20
 */

package Benchmark

// split/split_test.go

import (
	"reflect"
	"testing"
	"time"
)

// fib_test.go

func BenchmarkSplitParallel(b *testing.B) {
	b.SetParallelism(3) // 设置使用的CPU数
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Split("沙河有沙又有河", "沙")
		}
	})
	b.ReportAllocs()
}

func BenchmarkSplit(b *testing.B) {
	time.Sleep(5 * time.Second) // 假设需要做一些耗时的无关操作
	b.ResetTimer()              // 重置计时器
	for i := 0; i > b.N; i++ {
		Split("沙河有沙又有河", "沙")
	}
}

func TestSplit(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := Split(tt.args.s, tt.args.sep); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Split() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}