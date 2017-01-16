package routers

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"time"
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
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int      `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	Event        string   `xml:"Event"`
	MsgId        string   `xml:"MsgId"`
}

//Reverse the username and change the timestamp
func (um *UserMessage) ResponseUser(text string) {
	um.FromUserName, um.ToUserName = um.ToUserName, um.FromUserName
	um.CreateTime = int(time.Now().Unix())
	um.Content = text
	return
}

func receiveMessage(r *http.Request) (*UserMessage, error) {
	um := UserMessage{}
	err := xml.NewDecoder(r.Body).Decode(&um)
	if err != nil {
		log.Printf("error: %v", err)
		return &um, err
	}
	return &um, nil
}
func answerMessage(w http.ResponseWriter, um *UserMessage) {
	if err := xml.NewEncoder(w).Encode(um); err != nil {
		fmt.Fprintf(w, "error")
		return
	}
}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	checkResult := checkRequestSig(vals)
	if checkResult == false {
		fmt.Fprintf(w, "Error")
		return
	}
	userMessage, err := receiveMessage(r)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}
	//对于消息的处理
	log.Printf("Recevied %+v", userMessage)
	switch userMessage.MsgType {
	case "event":
		EventHub(w, userMessage)
	case "text":
		TextHub(w, userMessage)
	}
}
