package main

/**
 * Created by @CaomaoBoy on 2021/1/2.
 *  email:<115882934@qq.com>
 */
import (
	"github.com/gomodule/redigo/redis"
	"time"
)


var RedisClient = &redis.Pool{
	//连接方法
	Dial:            func() (redis.Conn,error) {
	c,err := redis.Dial("tcp","127.0.0.1:6379")
	if err != nil {
	return nil,err
	}
	err=c.Send("auth","123456")
	c.Do("SELECT",0)

	return c,nil
},
	//DialContext:     nil,
	//TestOnBorrow:    nil,
	//最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
	MaxIdle:         1,
	//最大的激活连接数，表示同时最多有N个连接
	MaxActive:       10,
	//最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
	IdleTimeout:     180 * time.Second,
	//Wait:            false,
	//MaxConnLifetime: 0,
}


//type Redis struct {
//	RedisClient *redis.Pool
//}
