package main

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func gameInit(data Post) { //hangman游戏逻辑
	const NUM int = 26
	guesses := 6
	wordList := []string{"apiary", "beetle", "cereal", "danger", "ensign", "florid", "garage", "health", "insult",
		"jackal", "keeper", "loaner", "manage", "nonce", "onset", "plaid", "quilt", "remote",
		"stolid", "train", "useful", "valid", "whence", "xenon", "yearn", "zippy"}
	rand.Seed(time.Now().UnixNano())
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
	err := HttpGetMessage(URL, initReq, res)
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
		err := HttpGetMessage(URL, req, res)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func gameService(data Post) {
	info := FindInfo(data.UserID)
	if len(info.Target) > 0 {
		data.Message = strings.ToLower(data.Message)
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
			err := HttpGetMessage(URL, req, res)
			if err != nil {
				log.Println(err.Error())
				return
			}
			if info.Guesses < 1 {
				DelInfo(data.UserID)
				req = request{
					GroupID: data.GroupID,
					Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]游戏结束！",
					Token:   AccessToken,
				}
				err = HttpGetMessage(URL, req, res)
				if err != nil {
					log.Println(err.Error())
					return
				}
				log.Println(res)
				req := MuteRequest{
					GroupID:  data.GroupID,
					UserID:   data.UserID,
					Duration: 300,
					Token:    AccessToken,
				}
				err = HttpGetMute(URLMute, req, res)
				if err != nil {
					log.Println(err.Error())
					return
				}
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
				err := HttpGetMessage(URL, req, res)
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
			err := HttpGetMessage(URL, req, res)
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
				err := HttpGetMessage(URL, req, res)
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
			err := HttpGetMessage(URL, req, res)
			if err != nil {
				log.Println(err.Error())
				return
			}
			log.Println(res)
		}
	}
}

func HotPotatoInit(data Post) {
	//寻找是否有进行中的游戏
	ttl := TTLList(data.GroupID)
	log.Println(ttl)
	if ttl < 1 {
		//如果没有进行中的游戏，则开启新的游戏队列
		CreateList(data.UserID, data.GroupID)
		req := request{
			GroupID: data.GroupID,
			Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]你开启了烫手山芋游戏，需要至少两人开始游戏，目前玩家人数:1\n其他人可输入“报名”加入游戏(一分钟内人数不足游戏自动结束)",
			Token:   AccessToken,
		}
		var res string
		err := HttpGetMessage(URL, req, res)
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Println(res)
		time.Sleep(time.Second * 60)
		//等待六十秒后，看是否有玩家加入游戏
		list := FindList(data.GroupID)
		if len(list) == 1 {
			req := request{
				GroupID: data.GroupID,
				Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]人数不足！游戏结束！",
				Token:   AccessToken,
			}
			var res string
			err := HttpGetMessage(URL, req, res)
			if err != nil {
				log.Println(err.Error())
				return
			}
			log.Println(res)
			return
		}
	} else {
		req := request{
			GroupID: data.GroupID,
			Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]目前有一场游戏在进行中，请输入“报名”加入游戏!",
			Token:   AccessToken,
		}
		var res string
		err := HttpGetMessage(URL, req, res)
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Println(res)
	}
}

func HotPotatoService(data Post) {
	//检测是否存在开启的游戏队列
	ttl := TTLList(data.GroupID)
	if ttl > 0 {
		list := FindList(data.GroupID)
		//检测输入报名的玩家是否已经报名
		for _, v := range list {
			if v == data.UserID {
				req := request{
					GroupID: data.GroupID,
					Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]你已报名！",
					Token:   AccessToken,
				}
				var res string
				err := HttpGetMessage(URL, req, res)
				if err != nil {
					log.Println(err.Error())
					return
				}
				log.Println(res)
				return
			}
		}
		//如果没有报名则插入队列中
		InsertToList(data.UserID, data.GroupID)
		list = FindList(data.GroupID)
		playerNum := len(list)
		req := request{
			GroupID: data.GroupID,
			Message: "[CQ:at,qq=" + strconv.FormatInt(data.UserID, 10) + "]你加入了烫手山芋游戏，目前玩家人数:" + strconv.Itoa(playerNum) +
				"\n其他人可输入“报名”加入游戏" +
				"\n无新的操作游戏将在五秒内开始",
			Token: AccessToken,
		}
		var res string
		err := HttpGetMessage(URL, req, res)
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Println(res)
		time.Sleep(time.Second * 5)
		//睡眠五秒后查看列表是否有玩家数量更改，没有更改游戏开始，更改了则等下一个线程执行任务
		list = FindList(data.GroupID)
		if len(list) == playerNum {
			DelList(data.GroupID)
			potatoes := []Potato{
				{
					Name:       "一级山芋",
					Time:       5,
					Percentage: 60,
				},
				{
					Name:       "二级山芋",
					Time:       10,
					Percentage: 50,
				},
				{
					Name:       "三级山芋",
					Time:       15,
					Percentage: 40,
				},
				{
					Name:       "四级山芋",
					Time:       20,
					Percentage: 30,
				},
				{
					Name:       "五级山芋",
					Time:       30,
					Percentage: 20,
				},
			}
			//随机生成山芋
			rand.Seed(time.Now().UnixNano())
			potato := potatoes[rand.Int()%5]
			req = request{
				GroupID: data.GroupID,
				Message: "山芋等级:" + potato.Name + "\n山芋爆炸概率:" + strconv.Itoa(potato.Percentage) + "%\n禁言时间：" + strconv.Itoa(potato.Time) + "min",
				Token:   AccessToken,
			}
			err = HttpGetMessage(URL, req, res)
			if err != nil {
				log.Println(err.Error())
				return
			}
			log.Println(res)
			//从第一个玩家开始遍历爆炸，爆炸游戏结束，玩家被禁言，未结束直到最后一个玩家，回到第一个玩家
			for i := 0; i < len(list); i++ {
				rand.Seed(time.Now().UnixNano())
				r := rand.Int() % 101
				if r <= potato.Percentage {
					req = request{
						GroupID: data.GroupID,
						Message: "山芋抛向了[CQ:at,qq=" + strconv.FormatInt(list[i], 10) + "]，然后爆炸了！",
						Token:   AccessToken,
					}
					err = HttpGetMessage(URL, req, res)
					if err != nil {
						log.Println(err.Error())
						return
					}
					log.Println(res)
					muteReq := MuteRequest{
						GroupID:  data.GroupID,
						UserID:   list[i],
						Duration: potato.Time * 60,
						Token:    AccessToken,
					}
					err = HttpGetMute(URLMute, muteReq, res)
					if err != nil {
						log.Println(err.Error())
						return
					}
					return
				} else {
					req = request{
						GroupID: data.GroupID,
						Message: "山芋抛向了[CQ:at,qq=" + strconv.FormatInt(list[i], 10) + "]，但是并没有发生什么",
						Token:   AccessToken,
					}
					err = HttpGetMessage(URL, req, res)
					if err != nil {
						log.Println(err.Error())
						return
					}
					log.Println(res)
					if i == len(list)-1 {
						i = 0
					}
				}
				time.Sleep(time.Millisecond * 800)
			}
		}
	}
}
