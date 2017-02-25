package main

//http://www.darkreading.com/
import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go"
	"github.com/zhangmingkai4315/weichat-golang-backend/utils"
)

//https://threatpost.com/blog/
const (
	DarkReadingURL   = "http://www.darkreading.com/"
	DarkReadingTable = "darkreading"
)

type DarkReadingItem struct {
	Title        string
	Link         string
	PostTime     time.Time
	User         string
	UserLink     string
	ShortContent string
	CoverImg     string
	Md5Sum       string
}

type DarkReadingList struct {
	DRL          []DarkReadingItem
	URL          string
	QueryStatus  utils.Status
	ErrorMessage string
}

func NewDarkReadingItem() *DarkReadingItem {
	return &DarkReadingItem{}
}
func NewDarkReadingList(url string) *DarkReadingList {
	drl := []DarkReadingItem{}
	return &DarkReadingList{DRL: drl, QueryStatus: utils.QUERY_INIT, URL: url}
}

func (self DarkReadingItem) String() string {
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

func darkreadingParse(doc *goquery.Document, drl *DarkReadingList) error {
	log.Printf("Start parse data")
	doclist := doc.Find("header.strong")
	log.Println(doclist.Html())
	currentItem := NewDarkReadingItem()
	re := regexp.MustCompile("[0-9/]+")
	fmt.Printf("%q\n", re.FindString("hello,02/12/2017,ddwdc"))
	doclist.Each(func(i int, s *goquery.Selection) {
		if i > 10 {
			return
		}
		currentItem.Link, _ = s.Find("a").Attr("href")
		currentItem.Title = strings.TrimSpace(s.Find("a").Text())
		contentPoint := s.Next().Next().Next()
		currentItem.ShortContent = contentPoint.Find(".black.smaller").Text()
		currentItem.User = contentPoint.Find(".allcaps.smaller").Text()
		contentPoint.Find(".allcaps.smaller").Next().Html()
		str := re.FindString(contentPoint.Find(".darkgray.smaller").Text())
		currentItem.PostTime, _ = utils.GetDateFromMMDDYYYY(str)

		currentItem.Md5Sum = utils.GetMD5Hash(currentItem.Title, currentItem.User)
		drl.DRL = append(drl.DRL, *currentItem)
		// log.Println("Current:" + currentItem.String())
		return
	})

	return nil
}
func (self *DarkReadingList) Start() error {

	log.Printf("Start send http request to %s", self.URL)
	self.QueryStatus = utils.QUERY_READY
	res, err := http.Get(self.URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	utfBody, err := iconv.NewReader(res.Body, "iso-8859-1", "utf-8")
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return err
	}
	log.Printf("Receive data from %s", self.URL)
	// 进行解析
	self.QueryStatus = utils.QUERY_RUNNING
	if err := darkreadingParse(doc, self); err != nil {
		return err
	}
	self.QueryStatus = utils.QUERY_STOPPED
	return nil
}

func (self *DarkReadingList) Save() error {
	db, err := utils.NewDatabase()
	defer db.Close()
	if err != nil {
		return err
	}
	log.Println("Begin save data to database...")
	for _, row := range self.DRL {
		_, err := db.Exec("INSERT INTO darkreading(title, link, post_date,user_name, short_content,md5) "+
			"VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE md5 = ?",
			row.Title,
			row.Link,
			row.PostTime,
			row.User,
			row.ShortContent,
			row.Md5Sum,
			row.Md5Sum)
		if err != nil {
			log.Printf("Saving To Database/Hackernews Fail:%f", err.Error())
		}
	}
	return nil

}

// ShowList will print all the list inside of threadpost list
func (self *DarkReadingList) ShowList() {
	for _, item := range self.DRL {
		fmt.Println(item.String())
	}
}
func (self *DarkReadingList) ShowStatus() utils.Status {
	return self.QueryStatus
}
func main() {
	drl := NewDarkReadingList(DarkReadingURL)
	err := drl.Start()
	if err != nil {
		log.Printf("Processing DarkReading Error %s", err.Error())
		return
	}
	if drl.ShowStatus() == utils.QUERY_ERROR {
		log.Printf("Query DarkReading error:%s", drl.ErrorMessage)
		return
	}
	drl.ShowList()
	if drl.ShowStatus() == utils.QUERY_STOPPED {
		err := drl.Save()
		if err != nil {
			log.Printf("Saving DarkReading Error %s", err.Error())
		}
		return
	}
}
