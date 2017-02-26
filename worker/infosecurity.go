package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/zhangmingkai4315/weichat-golang-backend/utils"
	"log"
	"time"
	//"strings"
	"strings"
)

const (
	InfoSecurityURL   = "https://www.infosecurity-magazine.com/news/"
	InfoSecurityTable = "infosecurity"
)

type InfoSecurityItem struct {
	Title        string
	Link         string
	PostTime     time.Time
	ShortContent string
	CoverImg     string
	Md5Sum       string
}

type InfoSecurityList struct {
	ISL          []InfoSecurityItem
	URL          string
	QueryStatus  utils.Status
	ErrorMessage string
}

func NewInfoSecurityItem() *InfoSecurityItem {
	return &InfoSecurityItem{}
}
func NewInfoSecurityList(url string) *InfoSecurityList {
	tpl := []InfoSecurityItem{}
	return &InfoSecurityList{ISL: tpl, QueryStatus: utils.QUERY_INIT, URL: url}
}
func (self InfoSecurityItem) String() string {
	return fmt.Sprintf(
		"Title:%s\n"+
			"PostDate:%s\n"+
			"Link:%s\n"+
			"Image:%s\n"+
			"Content:%s\n"+
			"Md5:%s\n", self.Title, self.PostTime, self.Link, self.CoverImg, self.ShortContent, self.Md5Sum)
}
func (self *InfoSecurityList) Start() error {

	log.Printf("Start send http request to %s", self.URL)
	self.QueryStatus = utils.QUERY_READY
	doc, err := goquery.NewDocument(self.URL)
	if err != nil {
		return err
	}
	log.Printf("Receive data from %s", self.URL)
	// 进行解析
	self.QueryStatus = utils.QUERY_RUNNING
	if err := infosecParse(doc, self); err != nil {
		return err
	}
	self.QueryStatus = utils.QUERY_STOPPED
	return nil
}

func infosecParse(doc *goquery.Document, isl *InfoSecurityList) error {
	log.Printf("Start parse data")
	doclist := doc.Find(".webpage-item.with-thumbnail")
	currentItem := NewInfoSecurityItem()
	doclist.Each(func(i int, s *goquery.Selection) {
		//log.Println(s.Html())
		currentItem.CoverImg, _ = s.Find(".thumbnail").Attr("src")
		currentItem.CoverImg = "http:" + currentItem.CoverImg
		currentItem.Title = strings.TrimSpace(s.Find("h3").Text())
		currentItem.Link, _ = s.Find("a").First().Attr("href")
		currentItem.ShortContent = strings.TrimSpace(s.Find("p").First().Text())
		postTime, err := utils.GetDateForInfoSec((s.Find(".pub-date").Text()))
		if err == nil {
			currentItem.PostTime = postTime
		} else {
			log.Println(err.Error())
		}
		currentItem.Md5Sum = utils.GetMD5Hash(currentItem.Title)
		isl.ISL = append(isl.ISL, *currentItem)
		return
	})
	return nil
}

func (self *InfoSecurityList) Save() error {
	db, err := utils.NewDatabase()
	defer db.Close()
	if err != nil {
		return err
	}
	log.Println("Begin save data to database...")
	for _, row := range self.ISL {
		_, err := db.Exec("INSERT INTO infosecurity(title, link, post_date, short_content,cover_img,md5) "+
			"VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE md5 = ?",
			row.Title,
			row.Link,
			row.PostTime,
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

func (self *InfoSecurityList) ShowList() {
	for _, item := range self.ISL {
		fmt.Println(item.String())
	}
}

func (self *InfoSecurityList) ShowStatus() utils.Status {
	return self.QueryStatus
}

func main() {
	isl := NewInfoSecurityList(InfoSecurityURL)
	err := isl.Start()
	if err != nil {
		log.Printf("Processing ThreadPost Error %s", err.Error())
		return
	}
	if isl.ShowStatus() == utils.QUERY_ERROR {
		log.Printf("Query ThreadPost error:%s", isl.ErrorMessage)
		return
	}
	isl.ShowList()
	if isl.ShowStatus() == utils.QUERY_STOPPED {
		err := isl.Save()
		if err != nil {
			log.Printf("Saving ThreadPost Error %s", err.Error())
		}
		return
	}
}
