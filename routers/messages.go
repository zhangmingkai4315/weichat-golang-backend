package routers

import (
	"net/http"
	"fmt"
	"encoding/xml"
	"log"
)
//<xml>
//<ToUserName><![CDATA[toUser]]></ToUserName>
//<FromUserName><![CDATA[fromUser]]></FromUserName>
//<CreateTime>1348831860</CreateTime>
//<MsgType><![CDATA[text]]></MsgType>
//<Content><![CDATA[this is a test]]></Content>
//<MsgId>1234567890123456</MsgId>
//</xml>

type UserMessage struct {
	XML xml.Name `xml:"xml"`
	ToUserName    string   `xml:"ToUserName"`
	FromUserName    string   `xml:"FromUserName"`
	CreateTime int  `xml:"CreateTime"`
	MsgType  string   `xml:"MsgType"`
	Content string `xml:"Content"`
	MsgId string   `xml:"MsgId"`
}

func receiveMessage(r *http.Request)(*UserMessage,error){
	um := UserMessage{}
	err:=xml.NewDecoder(r.Body).Decode(&um)
	if err != nil {
		log.Printf("error: %v", err)
		return &um,err
	}
	return &um,nil
}

func PostMessage (w http.ResponseWriter, r *http.Request){
	vals := r.URL.Query()
	checkResult := checkRequestSig(vals)
	if checkResult == false {
		fmt.Fprintf(w, "Error")
		return
	}
	userMessage,err:=receiveMessage(r)
	if err!=nil{
		fmt.Fprintf(w,"Error")
		return
	}
	//对于消息的处理
	log.Printf("Recevied %+v",userMessage)
	fmt.Fprintf(w,"Recevied")
}
