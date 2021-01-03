package main

import (
	"bytes"
	"os"
	"runtime"
	"strconv"
	"unsafe"
)

/**
 * Created by @CaomaoBoy on 2021/1/2.
 *  email:<115882934@qq.com>
 */

func Try(f func(),Handler func(interface{}),group *WatiGroup){
	defer func() {
		if err := recover();err != nil{
			err := PanicTrace(10)
			Handler(string(err))
		}
	}()
	f()
	if group != nil{
		group.Done(1)
	}
}
func Itoa(num interface{}) string {
	switch n := num.(type) {
	case int8:
		return strconv.FormatInt(int64(n), 10)
	case int16:
		return strconv.FormatInt(int64(n), 10)
	case int32:
		return strconv.FormatInt(int64(n), 10)
	case int:
		return strconv.FormatInt(int64(n), 10)
	case int64:
		return strconv.FormatInt(int64(n), 10)
	case uint8:
		return strconv.FormatUint(uint64(n), 10)
	case uint16:
		return strconv.FormatUint(uint64(n), 10)
	case uint32:
		return strconv.FormatUint(uint64(n), 10)
	case uint:
		return strconv.FormatUint(uint64(n), 10)
	case uint64:
		return strconv.FormatUint(uint64(n), 10)
	case []byte:
		return  *(*string)(unsafe.Pointer(&n))
	case bool:
		if n == false{return "0"}else{return "1"}
	case float32:
		return strconv.FormatFloat(float64(n), 'E', -1, 32)
	case float64:
		return strconv.FormatFloat(n, 'E', -1, 64)//float64
	case string:
		return n
	}
	return ""
}

// PanicTrace trace panic stack info.
func PanicTrace(kb int) []byte {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return stack
}



func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}