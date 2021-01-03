package main

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

/**
 * Created by @CaomaoBoy on 2021/1/2.
 *  email:<115882934@qq.com>
 */
type Redis_Ser interface {
	GeyHsetKey() string
}

type AXCC struct {
	B string
}


type Asd struct {
	Id int `cborm:"pk"`
	A *AXCC `cborm:"name:a " cc:"asdasd"`
	Count int `cborm:"name:ct" `
	Num float64 `cborm:"name:nm"`
	SliceTest []int `cborm:"name:st"`
	MapTest map[int]string `cborm:"name:mt"`
}

func (sx *Asd) GeyHsetKey() string{
	return "cmb:asd"
}


func GetAll(f Redis_Ser) ([]Redis_Ser,error){
	return nil,nil
}

func Get(f func() Redis_Ser) (Redis_Ser,error){
	val := f()
	conn := RedisClient.Get()
	defer RedisClient.Close()
	//反射解析
	t := reflect.TypeOf(val)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("Check type error not Struct")
	}
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	fieldNum := t.NumField()
	var Pk string
	for i := 0; i < fieldNum; i++ {
		if Pk != ""{
			break
		}
		tags := strings.Split(string(t.Field(i).Tag), "\"")
		for j:=0;j<len(tags) -1;j+=2{
			if strings.TrimSpace(tags[j]) == "cborm:"{
				kvs := strings.Split(tags[j+1]," ")
				if len(kvs) == 1 && kvs[0] == "pk"{
					Pk =  Itoa(v.Field(i).Interface())
					Pk = val.GeyHsetKey() + ":" + Pk
					break
				}

			}
		}
	}

	if tval, err := redis.Values(conn.Do("HGETALL",Pk));err != nil{
		return nil,err
	}else{
		//获取结构体的
		for i := 0; i < fieldNum; i++ {
			fieldInfo := v.Type().Field(i) // a reflect.StructField
			tag := fieldInfo.Tag           // a reflect.StructTag
			name := tag.Get("cborm")
			if name == "" {
				name = strings.ToLower(fieldInfo.Name)
			}
			//去掉逗号后面内容 如 `json:"voucher_usage,omitempty"`
			index := 0
			kvs := strings.Split(name, " ")
			namekv := kvs[index]
			for !strings.Contains(namekv,"name") && index  < len(kvs) -1 {
				index +=1
				namekv = kvs[index]
			}
			if strings.Contains(namekv,"name"){
				name = strings.Split(namekv, ":")[1]
				for j:=0;j<len(tval);j+=2{
					v8 := tval[j+1].([]uint8)
					strv := *(*string)(unsafe.Pointer(&v8))
					if string(tval[j].([]uint8)) == name{
						if v.CanSet(){
							switch v.Field(i).Kind() {
							case reflect.String:
								val := reflect.ValueOf(strv)
								v.Field(i).Set(val)
							case reflect.Struct:
								if x1,err := AssertType(t.Field(i).Type.String(),v8);err!= nil{
									log.Error("Handler Slice Type Error",err)
								}else{
									v.Field(i).Set(reflect.ValueOf(x1))
								}
							case reflect.Ptr:
								if x1,err := AssertType(t.Field(i).Type.String(),v8);err!= nil{
									log.Error("Handler Slice Type Error",err)
								}else{
									v.Field(i).Elem().Set(reflect.ValueOf(x1))
								}
							case reflect.Float32:
								float32, _ := strconv.ParseFloat(strv, 32)
								v.Field(i).Set(reflect.ValueOf(float32))
							case reflect.Float64:
								float64, _ := strconv.ParseFloat(strv, 64)
								v.Field(i).Set(reflect.ValueOf(float64))
							case reflect.Int32:
								int32, _ := strconv.ParseInt(strv, 10, 32)
								v.Field(i).Set(reflect.ValueOf(int32))
							case reflect.Int64:
								int64, _ := strconv.ParseInt(strv, 10, 64)
								v.Field(i).Set(reflect.ValueOf(int64))
							case reflect.Int:
								int, _ := strconv.Atoi(strv)
								v.Field(i).Set(reflect.ValueOf(int))
							case reflect.Slice:
								if x1,err := AssertType(t.Field(i).Type.String(),v8);err!= nil{
									log.Error("Handler Slice Type Error",err)
								}else{
									v.Field(i).Set(reflect.ValueOf(x1))
								}
							case reflect.Map:
								if x1,err := AssertType(t.Field(i).Type.String(),v8);err!= nil{
									log.Error("Handler Map Error",err)
								}else{
									v.Field(i).Set(reflect.ValueOf(x1))
								}
							}
						}
					}
				}
			}

		}
	}

	return val,nil

}



func Save(f func() Redis_Ser) error{
	val := f()
	conn := RedisClient.Get()
	defer RedisClient.Close()
	//反射解析
	t := reflect.TypeOf(val)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("Check type error not Struct")
	}
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	fieldNum := t.NumField()
	rdis_map := RedisMap{
		Val: make([]interface{},0),
	}
	rdis_map.Val = append(rdis_map.Val,val.GeyHsetKey())
	var Pk string
	for i := 0; i < fieldNum; i++ {
		tags := strings.Split(string(t.Field(i).Tag), "\"")
		for j:=0;j<len(tags) -1;j+=2{
			if strings.TrimSpace(tags[j]) == "cborm:"{
				kvs := strings.Split(tags[j+1]," ")
				if len(kvs) == 1 && kvs[0] == "pk"{
					Pk =  Itoa(v.Field(i).Interface())
					rdis_map.Val[0] = rdis_map.Val[0].(string) + ":" + Pk
				}
				var Key string
				for _,vv := range kvs{
					kv := strings.Split(vv,":")
					if strings.TrimSpace(kv[0])  == "name"{
						Key = strings.TrimSpace(kv[1])
						rdis_map.Val = append(rdis_map.Val,Key)
					}
					switch v.Field(i).Kind(){
						case  reflect.String:
							rdis_map.Val = append(rdis_map.Val,string(v.Field(i).Bytes()))
						case  reflect.Struct:
							val,err := json.Marshal(v.Field(i).Interface())
							if err != nil{
								log.Error(err)
							}
							rdis_map.Val = append(rdis_map.Val,string(val))
						case reflect.Ptr:
							val,err := json.Marshal(v.Field(i).Elem().Interface())
							if err != nil{
								log.Error(err)
							}
							rdis_map.Val = append(rdis_map.Val,string(val))
						case reflect.Uint , reflect.Uint16, reflect.Uint32 , reflect.Uint64:
							s := strconv.Itoa(int(v.Field(i).Uint()))
							rdis_map.Val = append(rdis_map.Val,s)
						case reflect.Int ,reflect.Int16 , reflect.Int32, reflect.Int64 :
							s := strconv.Itoa(int(v.Field(i).Int()))
							rdis_map.Val = append(rdis_map.Val,s)
						case  reflect.Float64 , reflect.Float32 :
							s := Itoa(v.Field(i).Float())
							rdis_map.Val = append(rdis_map.Val,s)
						case reflect.Bool:
							s := Itoa(v.Field(i).Bool())
							rdis_map.Val = append(rdis_map.Val,s)
						case reflect.Map:
							val,err := json.Marshal(v.Field(i).Interface())
							if err != nil{
								log.Error(err)
							}
							rdis_map.Val = append(rdis_map.Val,string(val))
						case reflect.Slice:
							val,err := json.Marshal(v.Field(i).Interface())
							if err != nil{
								log.Error(err)
							}
							rdis_map.Val = append(rdis_map.Val,string(val))
						}
				}

			}
		}
	}
	args := redis.Args{}
	for _,v := range rdis_map.Val{
		args = args.Add(v)
	}
	if _, err := conn.Do("HMSET",args...);err != nil{
		return err
	}
	if _, err := conn.Do("SADD",val.GeyHsetKey(),Pk);err != nil{
		return err
	}
	return nil

}

type RedisMap struct {
	Val []interface{}
}