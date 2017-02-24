package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const POST_MENU_URL = "https://api.weixin.qq.com/cgi-bin/menu/addconditional?access_token="
const DELET_MENU_URL = "https://api.weixin.qq.com/cgi-bin/menu/delconditional?access_token="
const QUERY_MENU_URL = "https://api.weixin.qq.com/cgi-bin/menu/get?access_token="
const Get_Current_URL = "https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token="

type Menu struct {
	Token string
}

const postMenuJson = `{
"button":
        [
            {
                "type": "click",
                "name": "Bind漏洞",
                "key":  "bind"
            },
            {
                "name": "安全新闻",
                "sub_button":
                [
                    {
                        "type": "click",
                        "name": "DarkReader",
                        "key":"darkreader"
                    },
                    {
                        "type": "click",
                        "name": "Hackernews",
                        "key":"hackernews"
                    },
                    {
                        "type": "click",
                        "name": "ThreatPost",
                        "key":"threatpost"
                    },
                    {
                        "type": "click",
                        "name": "Infosec",
                        "key":"infosecurity"
                    },
                ]
            }
          ]
}`

//type MenuObject struct {
//	Button []struct {
//		TypeName   string `json:type`
//		Name       string `json:name`
//		Key        string `json:key`
//		Sub_button []struct {
//			TypeName string `json:type`
//			Name     string `json:name`
//			Url      string `json:key`
//			Key      string `json:key`
//		}
//	}
//	Matchrule struct {
//		TagId          string
//		Sex            string
//		Country        string
//		province       string
//		City           string
//		ClientPlatForm string
//		Language       string
//	}
//}

func (menu *Menu) Create() (string, error) {
	postUrl := POST_MENU_URL + menu.Token
	postMenuJsonDate := []byte(postMenuJson)
	resp, err := http.Post(postUrl, "application/json", bytes.NewBuffer(postMenuJsonDate))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (menu *Menu) Query() (string, error) {
	postUrl := QUERY_MENU_URL + menu.Token
	resp, err := http.Get(postUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (menu *Menu) Delete() (string, error) {
	postUrl := DELET_MENU_URL + menu.Token
	resp, err := http.Get(postUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (menu *Menu) GetCurrentMenu() (string, error) {
	postUrl := Get_Current_URL + menu.Token
	resp, err := http.Get(postUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
