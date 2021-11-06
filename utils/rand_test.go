package utils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRandRangeNum(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandRangeNum(tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("RandRangeNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIncrStr(t *testing.T) {
	type args struct {
		i      int
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{rand.Intn(100) * 10, 4},
		},
	}
	for _, tt := range tests {
		fmt.Println("-------------------", "测试开始：", "[i]:", tt.args.i, "[length]:", tt.args.length, "-----------------------")
		i := tt.args.i
		ch := time.After((time.Minute) / 2) //设定一个超时时间，当计时到期则退出循环
		for {
			got := IncrStr(i, tt.args.length)
			fmt.Println(got)
			i++
			select {
			case t := <-ch:
				fmt.Println(t.Format(time.RFC3339), "到时退出", "【i】", i)
				return
			default:
			}
			time.Sleep(time.Second)
		}
	}
}
