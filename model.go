package main

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
type Potato struct {
	Name       string //山芋名称
	Time       int    //爆炸后禁言时间
	Percentage int    //爆炸概率（百分数)
}
type MuteRequest struct {
	GroupID  int64 `json:"group_id"`
	UserID   int64 `json:"user_id"`
	Duration int   `json:"duration"`
	Token string `json:"access_token"`
}
