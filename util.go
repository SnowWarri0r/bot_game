package main

import (
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"strconv"
	"strings"
	"time"
)

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

func HttpGetMessage(url string, data request, result string) error {
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

func HttpGetMute(url string, data MuteRequest, result string) error {
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
	params.Set("user_id", strconv.FormatInt(data.UserID, 10))
	params.Set("duration", strconv.Itoa(data.Duration))
	params.Set("access_token", data.Token)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	log.Println(urlPath)
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
