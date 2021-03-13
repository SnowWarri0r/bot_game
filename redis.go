package main

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

var pool *redis.Pool

func ConnectRedis() {
	pool = &redis.Pool{
		MaxIdle:     100,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			//数据库配置
			return redis.Dial("tcp", "127.0.0.1:8090", redis.DialPassword("passwd"))
		},
	}
}
func Pool() *redis.Pool {
	return pool
}
func DisconnectRedis() {
	if err := pool.Close(); err != nil {
		log.Println(err.Error())
	}
}
