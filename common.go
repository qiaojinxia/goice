package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strings"
	"sync"
)

/**
 * Created by @CaomaoBoy on 2021/1/2.
 *  email:<115882934@qq.com>
 */

//将所有对象 抽象为服务 有 初始化 启动 关闭 和 错误处理方法
type Server interface {
	Close() error //关闭服务
	Start() //开始服务
	InitByConfig(config *Config) error //初始化服务
	HandlerErr(interface{}) //处理服务器出错
}

type ServerContainer struct {
	*Servers
	data Server
	Once sync.Once
}

func (m *ServerContainer) Close() error{
	return m.data.Close()
}

func (m *ServerContainer) Start() {
	m.Once.Do(
		func() {
			//反射解析
			t := reflect.TypeOf(m.data)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			if t.Kind() != reflect.Struct {
				panic("Check type error not Struct")
			}
			v := reflect.ValueOf(m.data)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			fieldNum := t.NumField()
			for i := 0; i < fieldNum; i++ {
				tagName := t.Field(i).Name
				tags := strings.Split(string(t.Field(i).Tag), "\"")
				if len(tags) > 1 {
					tagName = tags[1]
				}
				if tags[0] == "inject:"{
					rtype := m.ContainerNames[tagName]
					if rtype == nil{
						log.Warn(fmt.Sprintf("Can't Find Register Type %s To Inject value to Field %s",tagName,t.Field(i).Name))
						continue
					}
					sc := reflect.ValueOf(m.containers[rtype].data)
					if  !v.Field(i).CanSet() {
						log.Warn(fmt.Sprintf("Can't Inject value to Field %s",t.Field(i).Name))
						continue
					}
					//vt := reflect.TypeOf(m.containers[rtype].data).Elem()
					//newoby := reflect.New(vt)
					//newoby.Elem().Set(reflect.ValueOf(m.containers[rtype].data).Elem())

					v.Field(i).Set(sc)
				}else if tags[0] == "cfi:" && m.confinject{
					//配置文件反射 tag 格式 cfi:field
					val := m.Config.Read(t.String(),tagName)
					sc := reflect.ValueOf(val)
					if  !v.Field(i).CanSet() {
						log.Warn(fmt.Sprintf("Can't Inject value to Field %s",t.Field(i).Name))
						continue
					}
					v.Field(i).Set(sc)
				}

			}
			err := m.data.InitByConfig(m.Config)
			if err != nil{
				panic(err)
			}
			m.AddWait(1)
			go Try(func() {
				m.data.Start()
			},m.data.HandlerErr,m.WatiGroup)
		})
}



type Servers struct {
	containers map[reflect.Type]*ServerContainer
	lazyInit bool
	confinject bool
	ContainerNames map[string]reflect.Type
	Config *Config

	*WatiGroup
}

func(s *Servers) SetLazyInit(val bool){
	s.lazyInit = val
}

func(s *Servers) SetConfInject(val bool){
	s.confinject = val
}
func NewServers() *Servers {
	return &Servers{containers: make(map[reflect.Type]*ServerContainer),ContainerNames: make(map[string]reflect.Type),WatiGroup:NewWatiGroup()}
}

func (s *Servers) RegisterGroup(v Server)  {

}

func (s *Servers) Register(v Server)  {
	sc := &ServerContainer{
		data: v,
		Once: sync.Once{},
		Servers:s,
	}
	msgType := reflect.TypeOf(v)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		panic("message pointer required")
	}
	if s.containers == nil {
		s.containers = make(map[reflect.Type]*ServerContainer)
	}
	s.ContainerNames[msgType.String()] = msgType
	if _,ok := s.containers[msgType];ok{
		panic(fmt.Sprintf("struct %s alerdy exists",msgType.Name()))
	}else{
		s.containers[msgType] = sc
	}
}

//初始化所有配置
func (s *Servers) start(mtype reflect.Type) {
	Try(func() {
		for tp,v := range s.containers{
			if s.lazyInit && mtype != tp{
				continue
			}
			v.Start()
		}
	}, func(i interface{}) {
		log.Error(i)
	},nil)

}
func (s *Servers) Get(mt interface{}) interface{}{
	msgType := reflect.TypeOf(mt)
	s.start(msgType)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		panic("message pointer required")
		return nil
	}
	return s.containers[msgType].data
}
//关闭所有
func (s *Servers) Close()  {
	for _,v := range s.containers{
		err := v.Close()
		if err != nil{
			log.Error(err)
		}
	}
}

func (s *Servers) InitConfig(path string) {
	//配置文件初始化
	confFile:= path
	if !Exists(confFile) {
		log.Error("配置文件不存在，程序退出")
		os.Exit(0)
	}
	config := &Config{}
	config.InitConfig(confFile)
	s.Config = config
}

