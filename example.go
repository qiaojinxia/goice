package main

import (
	"fmt"
	"time"
)

/**
 * Created by @CaomaoBoy on 2021/1/2.
 *  email:<115882934@qq.com>
 */

type Mya struct {
	val string
	Ival string `cfi:"val"`
	aa *AXX
}

type AXX struct {
	xx string

}

func (m *Mya) Close() error {
	fmt.Println("Server Close")
	return nil
}

func (m *Mya) HandlerErr(err interface{}){
	fmt.Println(err)
}


func (m *Mya) Start() {
	panic("抛出错误测试")
}

func (m *Mya) InitByConfig(config *Config) error {
	fmt.Println("获取到配置文件",config.Read("config","sku_id"))

	m.val = "初始化配置Mya"
	m.aa = new(AXX)

	fmt.Println(m.val)
	return nil
}

type AAA struct {
	Val *Mya `inject:"*main.Mya"`
}

func (m *AAA) Close() error {
	fmt.Println("Server Close")
	return nil
}

func (m *AAA) HandlerErr(err interface{}){
	fmt.Println(err)
}

func (m *AAA) Start() {
		for{
			tk := time.NewTicker(time.Second * 3)
			select {
			case <- tk.C:
				fmt.Printf("定时器 北京时间 %s\n" ,time.Now().Format("2006-01-02 15:04:05"))
			}
		}
}

func (m *AAA) InitByConfig(config *Config) error {
	fmt.Println("自动注入:",m.Val)
	return nil
}
//注册登录器
func main(){
	servermanger := NewServers()
	//初始化 懒加载
	servermanger.SetLazyInit(false)
	//从配置文件 加载
	servermanger.SetConfInject(true)
	servermanger.InitConfig("./conf.ini")
	servermanger.Register(&Mya{})
	servermanger.Register(&AAA{})
	val := servermanger.Get(&Mya{}).(*Mya)
	fmt.Println("修改前",val)
	val.aa = &AXX{} //修改结构体地址
	fmt.Println("修改后",val)
	val1 := servermanger.Get(&AAA{}).(*AAA)
	fmt.Println("修改后Copy",val1.Val)
	servermanger.Wait()

}