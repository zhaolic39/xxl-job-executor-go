package task

import (
	"context"
	xxl "github.com/zhaolic39/xxl-job-executor-go"
)

func Panic(cxt context.Context, param xxl.RunReq) (msg string) {
	panic("test panic")
	return
}
