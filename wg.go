package main

import (
	"sync"
)

/**
 * Created by @CaomaoBoy on 2021/1/2.
 *  email:<115882934@qq.com>
 */

type WatiGroup struct {
	wg *sync.WaitGroup
	lk sync.RWMutex
	count int
}

func NewWatiGroup() *WatiGroup {
	return &WatiGroup{wg: &sync.WaitGroup{}}
}

func(wg *WatiGroup) AddWait(num int){
	wg.lk.Lock()
	defer wg.lk.Unlock()
	wg.count += num
	wg.wg.Add(num)
}


func(wg *WatiGroup) Done(num int){
	wg.lk.Lock()
	defer wg.lk.Unlock()
	wg.count -= num
	for i:=0;i<num;i++{
		wg.wg.Done()
	}
}

func(wg *WatiGroup) Wait(){
	wg.wg.Wait()
}