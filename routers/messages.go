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
	XML          xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int      `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
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
	if answer, err := xml.Marshal(um); err != nil {
		fmt.Fprintf(w, "error")
		return
	} else {
		fmt.Fprintf(w, string(answer))
		return
	}
}

func UsageText() string {
	return "欢迎订阅DNS安全公众号,我们将实时为您反馈关于DNS及其相关的安全问题。\n" +
		"您可以回复如下的信息进行操作\n" +
		"1. 回复\"历史\"　查看历史信息" +
		"2. 回复\"攻击\"　查看dns安全事件" +
		"3. 其他内容，将自动转化为留言发送给管理员" +
		"感谢您的使用"

}

//　PostMessage　is a http handler, it will receive the text from the user and give them feedback

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
	userMessage.ResponseUser(UsageText())

	answerMessage(w, userMessage)
}
