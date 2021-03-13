package main

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
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

func DelInfo(id int64) {
	conn := Pool().Get()
	defer conn.Close()
	_, err := conn.Do("DEL", id)
	if err != nil {
		log.Println(err.Error())
	}
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

func CreateList(userID int64, groupId int64) {
	conn := Pool().Get()
	defer conn.Close()
	ok, err := redis.Bool(conn.Do("EXISTS", "GROUP_"+strconv.FormatInt(groupId, 10)))
	if err != nil {
		log.Println(err.Error())
		return
	}
	if ok {
		_, err = conn.Do("DEL", "GROUP_"+strconv.FormatInt(groupId, 10))
		if err != nil {
			log.Println(err.Error())
			return
		}
		_, err = conn.Do("RPUSH", "GROUP_"+strconv.FormatInt(groupId, 10), userID)
		if err != nil {
			log.Println(err.Error())
			return
		}
		_, err = conn.Do("EXPIRE", "GROUP_"+strconv.FormatInt(groupId, 10), 62)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	_, err = conn.Do("RPUSH", "GROUP_"+strconv.FormatInt(groupId, 10), userID)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = conn.Do("EXPIRE", "GROUP_"+strconv.FormatInt(groupId, 10), 62)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func InsertToList(userID int64, groupID int64) {
	conn := Pool().Get()
	defer conn.Close()
	ok, err := redis.Bool(conn.Do("EXISTS", "GROUP_"+strconv.FormatInt(groupID, 10)))
	if err != nil {
		log.Println(err.Error())
		return
	}
	if ok {
		_, err = conn.Do("RPUSH", "GROUP_"+strconv.FormatInt(groupID, 10), userID)
		if err != nil {
			log.Println(err.Error())
			return
		}
		_, err = conn.Do("EXPIRE", "GROUP_"+strconv.FormatInt(groupID, 10), 10)
	}
}
func FindList(groupID int64) []int64 {
	conn := Pool().Get()
	defer conn.Close()
	ok, err := redis.Bool(conn.Do("EXISTS", "GROUP_"+strconv.FormatInt(groupID, 10)))
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	if ok {
		data, err := redis.Int64s(conn.Do("LRANGE", "GROUP_"+strconv.FormatInt(groupID, 10), 0, -1))
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		return data
	}
	return nil
}
func TTLList(groupID int64) int {
	conn := Pool().Get()
	defer conn.Close()
	ok, err := redis.Bool(conn.Do("EXISTS", "GROUP_"+strconv.FormatInt(groupID, 10)))
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	if ok {
		ttl, err := redis.Int(conn.Do("TTL", "GROUP_"+strconv.FormatInt(groupID, 10)))
		if err != nil {
			log.Println(err.Error())
			return 0
		}
		return ttl
	}
	return 0
}
func DelList(groupID int64) {
	conn := Pool().Get()
	defer conn.Close()
	_, err := conn.Do("DEL", "GROUP_"+strconv.FormatInt(groupID, 10))
	if err != nil {
		log.Println(err.Error())
		return
	}
}
