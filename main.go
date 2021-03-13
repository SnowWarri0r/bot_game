package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const URL string = "url"
const URLMute string = "url"
const AccessToken string = "token"

func main() {
	ConnectRedis()
	e := gin.Default()
	e.POST("/hangman", hangMan())
	e.POST("/hotPotato", hotPotato())
	e.GET("/timer", timer())
	srv := &http.Server{
		Addr:    ":8081",
		Handler: e,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("listen: %s\n", err)
	}
	defer DisconnectRedis()
}
func hotPotato() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dat Post
		err := c.BindJSON(&dat)
		if err != nil {
			log.Println(err.Error())
		}
		if dat.PostType == "message" && dat.MessageType == "group" {
			if dat.Message == "烫手山芋游戏" {
				log.Println(dat)
				go HotPotatoInit(dat)
			} else if dat.Message == "报名" {
				go HotPotatoService(dat)
			}
		}
	}
}
func hangMan() gin.HandlerFunc {
	/*
		TODO:通过Redis缓存玩家信息，每次访问先搜索redis是否存在该玩家数据（通过QQ号)，存在则读取并使用玩家输入的信息进行游戏逻辑，然后更新缓存
		TODO:如果开始新游戏则先搜索是否存在缓存，存在则删除然后重新随机新数据并创建
	*/
	return func(c *gin.Context) {
		var dat Post
		err := c.BindJSON(&dat)
		if err != nil {
			log.Println(err.Error())
		}
		if dat.PostType == "message" && dat.MessageType == "group" {
			if dat.Message == "hangman游戏" {
				log.Println(dat)
				go gameInit(dat)
			} else if len(dat.Message) < 2 {
				go gameService(dat)
			}
		}
	}
}
func timer() gin.HandlerFunc {
	return func(c *gin.Context) {
		Timer()
	}
}
