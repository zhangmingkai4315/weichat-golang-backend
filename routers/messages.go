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
	PicUrl       string   `xml:"PicUrl"`
	Content      string   `xml:"Content"`
	Event        string   `xml:"Event"`
	MsgId        string   `xml:"MsgId"`
	MediaId      string   `xml:MediaId`
}

//ResponseUser 发送响应消息给用户：将消息中的from 和 to user 进行对调，并添加时间戳及传递的内容
func (um *UserMessage) ResponseUser() {
	um.FromUserName, um.ToUserName = um.ToUserName, um.FromUserName
	um.CreateTime = int(time.Now().Unix())
	return
}

func (um *UserMessage) ResponseUserWithContent(text string) {
	um.ResponseUser()
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
		fmt.Fprintf(w, "")
		return
	}
}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	// 线上的时候开启，接收微信的验证
	// 开发环境中关闭以下信息

	checkResult := checkRequestSig(vals)
	if checkResult == false {
		fmt.Fprintf(w, "Hash sorted list is not Equal Signature")
		return
	}

	userMessage, err := receiveMessage(r)
	if err != nil {
		fmt.Fprintf(w, "")
		return
	}
	//对于消息的处理
	//假如服务器无法保证在五秒内处理回复，则必须回复“success”或者“”（空串），否则微信后台会发起三次重试。
	log.Printf("Recevied %+v", userMessage)
	switch userMessage.MsgType {
	case "event":
		EventHub(w, userMessage)
	case "text":
		TextHub(w, userMessage)
	case "image":
		ImageHub(w, userMessage)

	default:
		fmt.Fprintf(w, "success")
	}
}
