package routers

import (
	"crypto/sha1"
	"fmt"
	"github.com/zhangmingkai4315/weichat-golang-backend/config"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func checkRequestSig(vals url.Values) bool {
	signature := vals.Get("signature")
	timestamp := vals.Get("timestamp")
	nonce := vals.Get("nonce")
	token := config.ConfigObj.Security.Token
	tmpArry := []string{token, timestamp, nonce}
	sort.Strings(tmpArry)
	tmpString := strings.Join(tmpArry, "")

	tmphash := sha1.New()
	tmphash.Write([]byte(tmpString))
	tmpBytesSlice := tmphash.Sum(nil)
	result := fmt.Sprintf("%x", tmpBytesSlice)
	if result == signature {
		return true
	} else {
		return false
	}
}
func WeiChatValidate(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	echostr := vals.Get("echostr")
	log.Printf("Reiceive Query : %s", r.URL.String())
	if checkResult := checkRequestSig(vals); checkResult == true {
		fmt.Fprintf(w, echostr)
		return
	} else {
		fmt.Fprintf(w, "Error")
		return
	}
}
