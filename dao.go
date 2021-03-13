package main

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
)

var NotExists = errors.New("not exists")

func CreateInfo(data *Info, id int64) {
	conn := Pool().Get()
	defer conn.Close()
	ok, err := redis.Bool(conn.Do("EXISTS", id))
	if err != nil {
		log.Println(err.Error())
		return
	}
	if ok {
		_, err = conn.Do("DEL", id)
		if err != nil {
			log.Println(err.Error())
			return
		}
		_, err = conn.Do("HMSET", redis.Args{}.Add(id).AddFlat(data)...)
		if err != nil {
			log.Println(err.Error())
			return
		}
		_, err = conn.Do("EXPIRE", id, 62)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	_, err = conn.Do("HMSET", redis.Args{}.Add(id).AddFlat(data)...)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = conn.Do("EXPIRE", id, 62)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func TTLInfo(id int64) int {
	conn := Pool().Get()
	defer conn.Close()
	ok, err := redis.Bool(conn.Do("EXISTS", id))
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	if ok {
		ttl, err := redis.Int(conn.Do("TTL", id))
		if err != nil {
			log.Println(err.Error())
			return 0
		}
		return ttl
	}
	return 0
}
func FindInfo(id int64) Info {
	conn := Pool().Get()
	defer conn.Close()
	ok, err := redis.Bool(conn.Do("EXISTS", id))
	if err != nil {
		log.Println(err.Error())
		return Info{}
	}
	var inf Info
	if ok {
		v, err := redis.Values(conn.Do("HGETALL", id))
		if err != nil {
			return Info{}
		}
		err = redis.ScanStruct(v, &inf)
		if err != nil {
			return Info{}
		}
	}
	return inf
}

func DelInfo(id int64) {
	conn := Pool().Get()
	defer conn.Close()
	_, err := conn.Do("DEL", id)
	if err != nil {
		log.Println(err.Error())
	}
}
