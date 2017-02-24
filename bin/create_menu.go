package main

import (
	"github.com/zhangmingkai4315/weichat-golang-backend/handler"
	"log"
)

func main() {
	if handler.Token.Token == "" {
		log.Panicln("Token不存在，无法发送请求")
	}
	menuHandler := handler.Menu{Token: handler.Token.Token}
	result, err := menuHandler.Create()
	if err != nil {
		log.Panic(err.Error())
	}
	log.Println(result)
}
