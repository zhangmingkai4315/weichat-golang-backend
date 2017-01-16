package routers

import "net/http"

func UsageText() string {
	return "欢迎订阅DNS安全公众号,我们将实时为您反馈关于DNS及其相关的安全问题。\n" +
		"您可以回复如下的信息进行操作。\n" +
		"1. 回复\"历史\" 查看历史信息。\n" +
		"2. 回复\"攻击\" 查看dns安全事件。\n" +
		"3. 发送其他内容，将自动转化为留言传递给管理员。\n" +
		"感谢您的使用。"

}
func ThanksForYourMessage() string {
	return "您的留言已收到，感谢您的意见反馈。\n" +
		"谢谢使用。"
}

//　PostMessage　is a http handler, it will receive the text from the user and give them feedback

func TextHub(w http.ResponseWriter, um *UserMessage) {
	switch um.Content {
	case "历史":
		um.ResponseUser(UsageText())
	case "攻击":
		um.ResponseUser(UsageText())
	default:
		um.ResponseUser(ThanksForYourMessage())
	}
	answerMessage(w, um)
}

func EventHub(w http.ResponseWriter, um *UserMessage) {
	switch um.Event {
	case "subscribe":
		fallthrough
	default:
		um.ResponseUser(UsageText())
	}
	answerMessage(w, um)
}
