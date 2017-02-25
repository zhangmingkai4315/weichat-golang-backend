package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/satori/go.uuid"
	"github.com/zhangmingkai4315/weichat-golang-backend/utils"
	"log"
	"time"
)

const (
	HackerNewsURL    = "https://news.ycombinator.com/"
	HackerNewsTables = "hackernews"
)

type HackerNewsItem struct {
	Id       string    `json:id`
	Title    string    `json:title`
	Link     string    `json:link`
	Time     time.Time `json:timestamp`
	From     string    `json:from`
	Score    string    `json:score`
	User     string    `json:user`
	UserLink string    `json:user_profile`
	Md5Sum   string    `json:md5`
}

type HackerNewsList struct {
	HNL          []HackerNewsItem `json:news`
	URL          string           `json:homepage`
	QueryStatus  utils.Status     `json:status`
	ErrorMessage string
}

func NewHackerNewsItem() *HackerNewsItem {
	// Creating UUID Version 4
	var id string
	id = fmt.Sprintf("%s", uuid.NewV4())
	return &HackerNewsItem{Id: id}
}

func (self HackerNewsItem) String() string {
	return fmt.Sprintf("------%s-------\n"+
		"Title:%s\n"+
		"PostDate:%s\n"+
		"Score:%s\n"+
		"User:%s\n"+
		"Link:%s\n"+
		"Web:%s\n"+
		"Md5:%s\n", self.Id, self.Title, self.Time, self.Score, self.User, self.Link, self.From, self.Md5Sum)
}

func NewHackerNewsList(url string) *HackerNewsList {
	hnl := []HackerNewsItem{}
	return &HackerNewsList{HNL: hnl, QueryStatus: utils.QUERY_INIT, URL: HackerNewsURL}
}
func hacknewsParse(doc *goquery.Document, hnl *HackerNewsList) error {
	log.Printf("Start parse data")
	doclist := doc.Find(".itemlist tbody tr")
	number := doclist.Length() / 3
	var currentItem *HackerNewsItem
	doclist.Each(func(i int, s *goquery.Selection) {
		if i >= 3*number {
			return
		}
		switch i % 3 {
		case 0:
			currentItem = NewHackerNewsItem()
			currentItem.Title = s.Find(".storylink").Text()
			currentItem.Link, _ = s.Find(".storylink").Attr("href")
			currentItem.From = s.Find(".sitestr").First().Text()
		case 1:
			currentItem.Score = s.Find(".score").Text()
			currentItem.User = s.Find(".hnuser").Text()
			href, exist := s.Find(".hnuser").Attr("href")
			if exist {
				currentItem.UserLink = HackerNewsURL + href
			}

			postTime, err := utils.GetDateFromString(s.Find(".age").First().Text())
			if err == nil {
				currentItem.Time = postTime
			}
			currentItem.Md5Sum = utils.GetMD5Hash(currentItem.Title, currentItem.User)
			hnl.HNL = append(hnl.HNL, *currentItem)
		case 2:
			fallthrough
		default:
			return
		}
		return
	})
	return nil
}
func (self *HackerNewsList) Start() error {

	log.Printf("Start send http request to %s", self.URL)
	self.QueryStatus = utils.QUERY_READY
	doc, err := goquery.NewDocument(self.URL)
	if err != nil {
		return err
	}
	log.Printf("Receive data from %s", self.URL)
	// 进行解析
	self.QueryStatus = utils.QUERY_RUNNING
	if err := hacknewsParse(doc, self); err != nil {
		return err
	}
	self.QueryStatus = utils.QUERY_STOPPED
	return nil
}

func (self *HackerNewsList) Save() error {
	db, err := utils.NewDatabase()
	defer db.Close()
	if err != nil {
		return err
	}
	log.Println("Begin save data to database...")
	for _, row := range self.HNL {
		_, err := db.Exec("INSERT INTO hackernews(uuid, title, link,post_date,score,user_name,user_profile,md5) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE uuid = ?",
			row.Id, row.Title, row.Link, row.Time, row.Score, row.User, row.UserLink, row.Md5Sum, row.Id)
		if err != nil {
			log.Printf("Saving To Database/Hackernews Fail:%f", err.Error())
		}
	}
	return nil

}

func (self *HackerNewsList) ShowList() {
	for _, item := range self.HNL {
		fmt.Println(item.String())
	}
}
func (self *HackerNewsList) ShowStatus() utils.Status {
	return self.QueryStatus
}

func main() {
	hnl := NewHackerNewsList(HackerNewsURL)
	err := hnl.Start()
	if err != nil {
		log.Printf("Processing HackerNews Error %s", err.Error())
		return
	}
	if hnl.ShowStatus() == utils.QUERY_ERROR {
		log.Printf("Query Hackernews error:%s", hnl.ErrorMessage)
		return
	}
	if hnl.ShowStatus() == utils.QUERY_STOPPED {
		err := hnl.Save()
		if err != nil {
			log.Printf("Saving HackerNews Error %s", err.Error())
		}
		return
	}
}
