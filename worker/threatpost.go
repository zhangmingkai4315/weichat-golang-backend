package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/zhangmingkai4315/weichat-golang-backend/utils"
)

//https://threatpost.com/blog/
const (
	ThreatpostURL    = "https://threatpost.com/blog/"
	ThreatpostTables = "threatpost"
)

type ThreatpostItem struct {
	Title        string
	Link         string
	PostTime     time.Time
	User         string
	UserLink     string
	ShortContent string
	CoverImg     string
	Md5Sum       string
}

type ThreatpostList struct {
	TPL          []ThreatpostItem
	URL          string
	QueryStatus  utils.Status
	ErrorMessage string
}

func NewThreadpostItem() *ThreatpostItem {
	return &ThreatpostItem{}
}
func NewThreadpostList(url string) *ThreatpostList {
	tpl := []ThreatpostItem{}
	return &ThreatpostList{TPL: tpl, QueryStatus: utils.QUERY_INIT, URL: url}
}

func (self ThreatpostItem) String() string {
	return fmt.Sprintf(
		"Title:%s\n"+
			"PostDate:%s\n"+
			"User:%s\n"+
			"UserLink:%s\n"+
			"Link:%s\n"+
			"Image:%s\n"+
			"Content:%s\n"+
			"Md5:%s\n", self.Title, self.PostTime, self.User, self.UserLink, self.Link, self.CoverImg, self.ShortContent, self.Md5Sum)
}

func threatpostParse(doc *goquery.Document, tpl *ThreatpostList) error {
	log.Printf("Start parse data")
	doclist := doc.Find("#latest-posts article")
	//log.Println(doc.Html())
	currentItem := NewThreadpostItem()
	doclist.Each(func(i int, s *goquery.Selection) {
		currentItem.CoverImg, _ = s.Find(".image-wrap a img").Attr("src")
		currentItem.Title = strings.TrimSpace(s.Find(".entry-title").Text())
		currentItem.Link, _ = s.Find(".read-more").Attr("href")
		currentItem.User = s.Find(".post-info .author a").Text()
		currentItem.UserLink, _ = s.Find(".post-info .author a").Attr("href")
		currentItem.ShortContent = strings.TrimSpace(s.Find(".categories").Next().Text())
		postTime, err := utils.GetDateForThreadPostLayout((s.Find("time").First().Text()))
		if err == nil {
			currentItem.PostTime = postTime
		} else {
			log.Println(err.Error())
		}
		currentItem.Md5Sum = utils.GetMD5Hash(currentItem.Title, currentItem.User)
		tpl.TPL = append(tpl.TPL, *currentItem)
		log.Println("Current:" + currentItem.String())
		return
	})
	return nil
}
func (self *ThreatpostList) Start() error {

	log.Printf("Start send http request to %s", self.URL)
	self.QueryStatus = utils.QUERY_READY
	doc, err := goquery.NewDocument(self.URL)
	if err != nil {
		return err
	}
	log.Printf("Receive data from %s", self.URL)
	// 进行解析
	self.QueryStatus = utils.QUERY_RUNNING
	if err := threatpostParse(doc, self); err != nil {
		return err
	}
	self.QueryStatus = utils.QUERY_STOPPED
	return nil
}

func (self *ThreatpostList) Save() error {
	db, err := utils.NewDatabase()
	defer db.Close()
	if err != nil {
		return err
	}
	log.Println("Begin save data to database...")
	for _, row := range self.TPL {
		_, err := db.Exec("INSERT INTO threadpost(title, link, post_date,user_name,user_profile, short_content,cover_img,md5) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE md5 = ?",
			row.Title,
			row.Link,
			row.PostTime,
			row.User,
			row.UserLink,
			row.ShortContent,
			row.CoverImg,
			row.Md5Sum,
			row.Md5Sum)
		if err != nil {
			log.Printf("Saving To Database/Hackernews Fail:%f", err.Error())
		}
	}
	return nil

}

func (self *ThreatpostList) ShowList() {
	for _, item := range self.TPL {
		fmt.Println(item.String())
	}
}
func (self *ThreatpostList) ShowStatus() utils.Status {
	return self.QueryStatus
}
func main() {
	tpl := NewThreadpostList(ThreatpostURL)
	err := tpl.Start()
	if err != nil {
		log.Printf("Processing ThreadPost Error %s", err.Error())
		return
	}
	if tpl.ShowStatus() == utils.QUERY_ERROR {
		log.Printf("Query ThreadPost error:%s", tpl.ErrorMessage)
		return
	}
	tpl.ShowList()
	if tpl.ShowStatus() == utils.QUERY_STOPPED {
		err := tpl.Save()
		if err != nil {
			log.Printf("Saving ThreadPost Error %s", err.Error())
		}
		return
	}
}
