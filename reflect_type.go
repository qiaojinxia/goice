package main

import (
	"encoding/json"
	"errors"
)

/**
 * Created by @CaomaoBoy on 2021/1/3.
 *  email:<115882934@qq.com>
 */

//类型断言
func AssertType(key string,sdata []byte) (interface{},error){
	switch key{
	case "[]int":
		v := make([]int,0)
		err := json.Unmarshal(sdata,&v)
		if err != nil{
			return nil,err
		}
		return v,nil
	case "map[int]string":
		v := make(map[int]string)
		err := json.Unmarshal(sdata,&v)
		if err != nil{
			return nil,err
		}
		return v,nil
	case "*main.AXCC":
		v :=  new(AXCC)
		err := json.Unmarshal(sdata,&v)
		if err != nil{
			return nil,err
		}
		return v,nil
	case "main.AXCC":
		v := AXCC{}
		err := json.Unmarshal(sdata,&v)
		if err != nil{
			return nil,err
		}
		return v,nil
	default:
		return nil,errors.New("不支持的反序列化类型!")
	}


}
