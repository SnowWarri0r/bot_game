package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	url2 "net/url"
	"strconv"
	"strings"
	"time"
)

type request struct {
	GroupID int64  `json:"group_id"`
	Message string `json:"message"`
	Token   string `json:"access_token"`
}
type Info struct {
	Target   string `redis:"target"`
	Attempt  string `redis:"attempt"`
	Guesses  int    `redis:"guesses"`
	BadChars string `redis:"bad_chars"`
}
type Post struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	Message     string `json:"message"`
	GroupID     int64  `json:"group_id"`
	UserID      int64  `json:"user_id"`
}
//配置文件，还需转至redis.go修改数据库配置
const URL string = "url"
const AccessToken string = "token"

func main() {
	ConnectRedis()
	e := gin.Default()
	e.POST("/hangman", hangMan())
	srv := &http.Server{
		Addr:    ":8081",
		Handler: e,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("listen: %s\n", err)
	}
	defer DisconnectRedis()
}

func HttpGet(url string, data request, result string) error {
	var client = http.Client{
		Timeout: 10 * time.Second,
	}
	params := url2.Values{}
	Url, err := url2.Parse(url)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	params.Set("group_id", strconv.FormatInt(data.GroupID, 10))
	params.Set("message", data.Message)
	params.Set("access_token", data.Token)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	response, err := client.Get(urlPath)
	if err != nil {
		return err
	}
	var body []byte
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	result = string(body)
	return nil
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

func gameInit(data Post) { //hangman游戏逻辑
	const NUM int = 26
	guesses := 6
	wordList := []string{"apiary", "beetle", "cereal", "danger", "ensign", "florid", "garage", "health", "insult",
		"jackal", "keeper", "loaner", "manage", "nonce", "onset", "plaid", "quilt", "remote",
		"stolid", "train", "useful", "valid", "whence", "xenon", "yearn", "zippy"}
	target := wordList[rand.Int()%NUM]
	attempt := strings.Repeat("-", len(target))
	var badChars = ""
	info := Info{
		Attempt:  attempt,
		Target:   target,
		Guesses:  guesses,
		BadChars: badChars,
	}
	CreateInfo(&info, data.UserID)
	initReq := request{
		GroupID: data.GroupID,
		Message: "来猜猜单词吧，我这有个单词，共有" + strconv.Itoa(len(target)) + "个字母。\n你共有六次猜错的机会\n你要猜的单词:" + attempt + "\n输入一个字母吧!(一分钟无操作游戏自动结束)",
		Token:   AccessToken,
	}
	var res string
	err := HttpGet(URL, initReq, res)
	if err != nil {
		log.Println(err.Error())
		return
	}
	time.Sleep(time.Second * 60)
	ttl := TTLInfo(data.UserID)
	inf := FindInfo(data.UserID)
	if (inf.Guesses == info.Guesses || inf.Attempt == info.Attempt) && ttl < 5 {
		req := request{
			GroupID: data.GroupID,
			Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]一分钟内无操作，游戏自动结束。",
			Token:   AccessToken,
		}
		err := HttpGet(URL, req, res)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func gameService(data Post) {
	info := FindInfo(data.UserID)
	if len(info.Target) > 0 {
		data.Message=strings.ToLower(data.Message)
		loc := strings.Index(info.Target, data.Message)
		if loc < 0 { //猜测错误逻辑
			info.Guesses--
			info.BadChars += data.Message
			CreateInfo(&info, data.UserID)
			req := request{
				GroupID: data.GroupID,
				Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]猜测错误!你还有" + strconv.Itoa(info.Guesses) +
					"次机会!\n错误的字符合计:" + info.BadChars + "\n你要猜测的单词:" + info.Attempt,
				Token: AccessToken,
			}
			var res string
			err := HttpGet(URL, req, res)
			if err != nil {
				log.Println(err.Error())
				return
			}
			if info.Guesses < 1 {
				DelInfo(data.UserID)
				req := request{
					GroupID: data.GroupID,
					Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]游戏结束！",
					Token:   AccessToken,
				}
				var res string
				err := HttpGet(URL, req, res)
				if err != nil {
					log.Println(err.Error())
					return
				}
				log.Println(res)
				return
			}
			time.Sleep(time.Second * 60)
			ttl := TTLInfo(data.UserID)
			inf := FindInfo(data.UserID)
			if (inf.Guesses == info.Guesses || inf.Attempt == info.Attempt) && ttl < 5 {
				req := request{
					GroupID: data.GroupID,
					Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]一分钟内无操作，游戏自动结束。",
					Token:   AccessToken,
				}
				DelInfo(data.UserID)
				err := HttpGet(URL, req, res)
				if err != nil {
					log.Println(err.Error())
					return
				}
				log.Println(res)
			}
			return
		} else { //猜测正确逻辑
			for ; loc < len(info.Target); {
				info.Attempt = info.Attempt[:loc] + data.Message + info.Attempt[loc+1:]
				log.Println(info.Attempt)
				loc = subStrIndex(info.Target, data.Message, loc+1)
				if loc < 0 {
					break
				}
			}
		}
		if info.Attempt != info.Target { //判断是否已经猜测对单词
			CreateInfo(&info, data.UserID)
			req := request{
				GroupID: data.GroupID,
				Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]猜测正确!\n错误的字符合计:" + info.BadChars + "\n你要猜测的单词:" + info.Attempt,
				Token:   AccessToken,
			}
			var res string
			err := HttpGet(URL, req, res)
			if err != nil {
				log.Println(err.Error())
				return
			}
			log.Println(res)
			time.Sleep(time.Second * 60)
			ttl := TTLInfo(data.UserID)
			inf := FindInfo(data.UserID)
			if (inf.Guesses == info.Guesses || inf.Attempt == info.Attempt) && ttl < 5 {
				req := request{
					GroupID: data.GroupID,
					Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]一分钟内无操作，游戏自动结束。",
					Token:   AccessToken,
				}
				DelInfo(data.UserID)
				err := HttpGet(URL, req, res)
				if err != nil {
					log.Println(err.Error())
					return
				}
				log.Println(res)
			}
		} else {
			DelInfo(data.UserID)
			req := request{
				GroupID: data.GroupID,
				Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]猜测正确!\n你猜测的单词为:" + info.Target,
				Token:   AccessToken,
			}
			var res string
			err := HttpGet(URL, req, res)
			if err != nil {
				log.Println(err.Error())
				return
			}
			log.Println(res)
		}
	}
}

func subStrIndex(str string, ch string, offset int) int {
	if offset >= len(str) {
		return -1
	}
	str = str[offset:]
	index := strings.Index(str, ch)
	if index < 0 {
		return -1
	}
	return offset + index
}
