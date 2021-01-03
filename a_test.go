package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

/**
 * Created by @CaomaoBoy on 2021/1/2.
 *  email:<115882934@qq.com>
 */
func Test_Save(t *testing.T) {
	Save(func() Redis_Ser {
		b := &Asd{
			Id: 996,
			A:     new(AXCC),
			Count: 123123,
			Num: 123213.098998,
			SliceTest:make([]int,0,10),
			MapTest: make(map[int]string),
		}
		b.SliceTest = append(b.SliceTest, 123)
		b.SliceTest = append(b.SliceTest, 345)
		b.SliceTest = append(b.SliceTest, 345)
		b.SliceTest = append(b.SliceTest, 7775)
		b.A.B= "ASDSAD"
		b.MapTest[123] = "234234"
		b.MapTest[234234] = "xxx4234"
		return b
	})
}

func Test_Get(t *testing.T) {
	val,_ := Get(func() Redis_Ser {
		return  &Asd{
			Id:996,
		}
	})
	fmt.Println(val.(*Asd))
}

type User struct {
	Id int
	Name string
}

func (u User) it() {
	panic("implement me")
}

type xxx interface {
	it()
}

func Test_B(tx *testing.T) {
	u := &User{1, "test"}

	t := reflect.TypeOf(u).Elem()
	//n := reflect.New(t).Elem().Interface()

	// 改成下面这样是正常的，但是在项目中并不知道是哪个结构体，也许是 Order, Pay...
	n := reflect.New(t).Elem().Interface()

	r, _ := json.Marshal(u)
	json.Unmarshal(r, &n)

	fmt.Println(n)
}