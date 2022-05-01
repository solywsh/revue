package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type httpPost struct {
	Time        int    `json:"time"`
	SelfId      int    `json:"self_id"`
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	SubType     string `json:"sub_type"`
	MessageId   int    `json:"message_id"`
	UserId      int    `json:"user_id"`
	Message     string `json:"message"`
	RawMessage  string `json:"raw_message"`
	Font        int    `json:"font"`
	Sender      struct {
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		Age      int    `json:"age"`
	} `json:"sender"`
}

func listenFromCqhttp(c *gin.Context) {
	var form httpPost
	if c.ShouldBind(&form) == nil {
		//fmt.Printf("%#v", form)
		info, _ := json.Marshal(form)
		fmt.Println(string(info))
	}
}

func main() {
	router := gin.Default()
	router.POST("/listenFromCqhttp", listenFromCqhttp)
	router.Run(":5000")
}
